package albums

import (
	"context"
	"errors"
	"net/http"
)

var (
	NullAlbum        = Album{}
	ErrAlbumNotFound = errors.New("album not found")

	// Excludes non app created albums. Google Photos doesn't allow manage non created albums through the API.
	// https://developers.google.com/photos/library/guides/manage-albums#adding-items-to-album
	excludeNonAppCreatedData = &ListOptions{ExcludeNonAppCreatedData: true}
)

// AlbumsService represents a Google Photos client for albums management.
type AlbumsService interface {
	AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*Album, error)
	GetById(ctx context.Context, id string) (*Album, error)
	GetByTitle(ctx context.Context, title string) (*Album, error)
	List(ctx context.Context) ([]Album, error)
	Update(ctx context.Context, album Album, updateMask []Field) (*Album, error)
}

// repository represents the Google Photos API client.
type repository interface {
	BatchAddMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	BatchRemoveMediaItemsAll(albumId string, mediaItemIds []string, ctx context.Context) error
	Create(title string, ctx context.Context) (*Album, error)
	Get(id string, ctx context.Context) (*Album, error)
	ListAll(options *ListOptions, ctx context.Context) ([]Album, error)
	ListAllAsync(options *ListOptions, ctx context.Context) (<-chan Album, <-chan error)
	Patch(album Album, fieldMask []Field, ctx context.Context) (*Album, error)
}

// Cache represents a cache service to store albums.
type Cache interface {
	GetAlbum(ctx context.Context, title string) (Album, error)
	PutAlbum(ctx context.Context, album Album) error
	InvalidateAlbum(ctx context.Context, title string) error
	InvalidateAllAlbums(ctx context.Context) error
}

// CachedAlbumsService implements a Google Photos client with cached results.
type CachedAlbumsService struct {
	repo  repository
	cache Cache
}

// AddMediaItems adds multiple media items to the specified album.
func (s CachedAlbumsService) AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.repo.BatchAddMediaItemsAll(albumId, mediaItemIds, ctx)
}

// RemoveMediaItems removes multiple media items from the specified album.
func (s CachedAlbumsService) RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.repo.BatchRemoveMediaItemsAll(albumId, mediaItemIds, ctx)
}

// Create adds and caches a new album to the repo.
func (s CachedAlbumsService) Create(ctx context.Context, title string) (*Album, error) {
	albumPtr, err := s.repo.Create(title, ctx)
	if err != nil {
		return &NullAlbum, err
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// GetById fetches and caches an album from the repo by id. It doesn't use the cache to look for it.
func (s CachedAlbumsService) GetById(ctx context.Context, id string) (*Album, error) {
	albumPtr, err := s.repo.Get(id, ctx)
	if err != nil {
		return &NullAlbum, err
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// GetByTitle fetches and caches an album from the repo by title. It tries to find it in the cache, first.
func (s CachedAlbumsService) GetByTitle(ctx context.Context, title string) (*Album, error) {
	album, err := s.cache.GetAlbum(ctx, title)
	if err == nil {
		return &album, nil // album was found in the cache
	}

	albumsC, errorsC := s.repo.ListAllAsync(excludeNonAppCreatedData, ctx)
	for {
		select {
		case item, ok := <-albumsC:
			if !ok {
				return &NullAlbum, ErrAlbumNotFound // there aren't more albums, album not found
			}
			if item.Title == title {
				return &item, s.cache.PutAlbum(ctx, item) // found, cache it and return it
			}
			if err := s.cache.PutAlbum(ctx, item); err != nil {
				return &NullAlbum, err // error when caching, returns
			}
		case err := <-errorsC:
			return &NullAlbum, err
		}
	}
}

// List fetches and caches all the albums from the repo.
func (s CachedAlbumsService) List(ctx context.Context) ([]Album, error) {
	result := make([]Album, 0)
	if err := s.cache.InvalidateAllAlbums(ctx); err != nil {
		return result, err
	}
	albumsC, errorsC := s.repo.ListAllAsync(excludeNonAppCreatedData, ctx)
	for {
		select {
		case item, ok := <-albumsC:
			if !ok {
				return result, nil
			}
			result = append(result, item)
			if err := s.cache.PutAlbum(ctx, item); err != nil {

				return result, err
			}
		case err := <-errorsC:
			return result, err
		}
	}
}

// Update updates album fields. updateMask argument can be used to update only selected fields. Currently only id, title
// and coverPhotoMediaItemId are read.
func (s CachedAlbumsService) Update(ctx context.Context, album Album, updateMask []Field) (*Album, error) {
	if err := s.cache.InvalidateAlbum(ctx, album.Title); err != nil {
		return nil, err
	}
	albumPtr, err := s.repo.Patch(album, updateMask, ctx)
	if err != nil {
		return &NullAlbum, err
	}
	err = s.cache.PutAlbum(ctx, *albumPtr)
	return albumPtr, err
}

// NewCachedAlbumsService returns a client of CachedAlbumsService.
func NewCachedAlbumsService(authenticatedClient *http.Client, options ...Option) CachedAlbumsService {
	var repo repository = defaultRepo(authenticatedClient)
	var albumCache Cache = defaultCache()

	for _, o := range options {
		switch o.Name() {
		case optkeyRepo:
			repo = o.Value().(repository)
		case optkeyCache:
			albumCache = o.Value().(Cache)
		}
	}

	return CachedAlbumsService{
		repo:  repo,
		cache: albumCache,
	}
}

const (
	optkeyRepo  = "repository"
	optkeyCache = "cache"
)

// Option represents a configurable parameter.
type Option interface {
	Name() string
	Value() interface{}
}

type option struct {
	name  string
	value interface{}
}

func (o option) Name() string       { return o.name }
func (o option) Value() interface{} { return o.value }

// WithRepository configures the Google Photos repository.
func WithRepository(s repository) Option {
	return &option{
		name:  optkeyRepo,
		value: s,
	}
}

// WithCache configures the cache.
func WithCache(s Cache) Option {
	return &option{
		name:  optkeyCache,
		value: s,
	}
}
