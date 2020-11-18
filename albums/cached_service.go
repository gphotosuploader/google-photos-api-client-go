package albums

import (
	"context"
	"errors"
	"net/http"
)

type Album struct {
	ID                    string
	Title                 string
	ProductURL            string
	IsWriteable           bool
	MediaItemsCount       string
	CoverPhotoBaseURL     string
	CoverPhotoMediaItemID string
}

// AlbumsService represents a Google Photos client for albums management.
type AlbumsService interface {
	AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*Album, error)
	GetById(ctx context.Context, id string) (*Album, error)
	GetByTitle(ctx context.Context, title string) (*Album, error)
	List(ctx context.Context) ([]Album, error)
}

var (
	NullAlbum        = Album{}
	ErrAlbumNotFound = errors.New("album not found")
)

// Repository represents the repository to store albums.
type Repository interface {
	AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*Album, error)
	Get(ctx context.Context, albumId string) (*Album, error)
	ListAll(ctx context.Context) ([]Album, error)
	GetByTitle(ctx context.Context, title string) (*Album, error)
}

// CachedAlbumsService implements a Google Photos client with cached results.
type CachedAlbumsService struct {
	repo  Repository
	cache Cache
}

// AddMediaItems adds multiple media items to the specified album.
func (s CachedAlbumsService) AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.repo.AddManyItems(ctx, albumId, mediaItemIds)
}

// RemoveMediaItems removes multiple media items from the specified album.
func (s CachedAlbumsService) RemoveMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.repo.RemoveManyItems(ctx, albumId, mediaItemIds)
}

// Create adds and caches a new album to the repo.
func (s CachedAlbumsService) Create(ctx context.Context, title string) (*Album, error) {
	albumPtr, err := s.repo.Create(ctx, title)
	if err != nil {
		return &NullAlbum, err
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// GetById fetches and caches an album from the repo by id. It doesn't use the cache to look for it.
func (s CachedAlbumsService) GetById(ctx context.Context, albumId string) (*Album, error) {
	albumPtr, err := s.repo.Get(ctx, albumId)
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
	a, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return &NullAlbum, err
	}
	return a, s.cache.PutAlbum(ctx, *a)
}

// List fetches and caches all the albums from the repo.
func (s CachedAlbumsService) List(ctx context.Context) ([]Album, error) {
	nullAlbums := make([]Album, 0)
	if err := s.cache.InvalidateAllAlbums(ctx); err != nil {
		return nullAlbums, err
	}
	albums, err := s.repo.ListAll(ctx)
	if err != nil {
		return nullAlbums, err
	}
	return albums, s.cache.PutManyAlbums(ctx, albums)
}

func defaultRepo(authenticatedClient *http.Client) Repository {
	return NewDuffplAlbumRepository(authenticatedClient)
}

func defaultCache() Cache {
	return NewCachitaCache()
}

// NewCachedAlbumsService returns a client of CachedAlbumsService.
func NewCachedAlbumsService(authenticatedClient *http.Client, options ...Option) CachedAlbumsService {
	var repo Repository = defaultRepo(authenticatedClient)
	var albumCache Cache = defaultCache()

	for _, o := range options {
		switch o.Name() {
		case optkeyRepo:
			repo = o.Value().(Repository)
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
func WithRepository(s Repository) Option {
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
