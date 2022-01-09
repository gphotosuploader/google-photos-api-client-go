package albums

import (
	"context"
	"errors"
	"net/http"
)

// Repository represents an album repository.
type Repository interface {
	AddManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	RemoveManyItems(ctx context.Context, albumId string, mediaItemIds []string) error
	Create(ctx context.Context, title string) (*Album, error)
	Get(ctx context.Context, albumId string) (*Album, error)
	ListAll(ctx context.Context) ([]Album, error)
	GetByTitle(ctx context.Context, title string) (*Album, error)
}

// Cache is used to store and retrieve previously obtained objects.
type Cache interface {
	// GetAlbum returns Album data from the cache corresponding to the specified title.
	// It will return ErrCacheMiss if there is no cached Album.
	GetAlbum(ctx context.Context, title string) (Album, error)

	// PutAlbum stores the Album data in the cache using the title as key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, GetAlbum, results in the original data.
	PutAlbum(ctx context.Context, album Album) error

	// PutManyAlbums stores many Album data in the cache using the title as key.
	PutManyAlbums(ctx context.Context, albums []Album) error

	// InvalidateAlbum removes the Album data from the cache corresponding to the specified title.
	// If there's no such Album in the cache, it will return nil.
	InvalidateAlbum(ctx context.Context, title string) error

	// InvalidateAllAlbums removes all key corresponding to albums
	InvalidateAllAlbums(ctx context.Context) error
}

// CachedAlbumsService implements a albums Google Photos client with cached results.
type CachedAlbumsService struct {
	repo  Repository
	cache Cache
}

var (
	// NullAlbum is a zero value Album.
	NullAlbum = Album{}

	ErrAlbumNotFound = errors.New("album not found")
)

// AddMediaItems adds multiple media item(s) to the specified album.
func (s CachedAlbumsService) AddMediaItems(ctx context.Context, albumId string, mediaItemIds []string) error {
	return s.repo.AddManyItems(ctx, albumId, mediaItemIds)
}

// RemoveMediaItems removes multiple media item(s) from the specified album.
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

// GetById fetches and caches an album from the repo by id.
// It does not use the cache to look for it.
// Returns ErrAlbumNotFound if the album does not exist.
func (s CachedAlbumsService) GetById(ctx context.Context, albumId string) (*Album, error) {
	albumPtr, err := s.repo.Get(ctx, albumId)
	if err != nil {
		return &NullAlbum, ErrAlbumNotFound
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// GetByTitle fetches and caches an album from the repo by title.
// It tries to find it in the cache, first.
// Returns ErrAlbumNotFound if the album does not exist.
func (s CachedAlbumsService) GetByTitle(ctx context.Context, title string) (*Album, error) {
	album, err := s.cache.GetAlbum(ctx, title)
	if err == nil {
		return &album, nil // album was found in the cache
	}
	a, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return &NullAlbum, ErrAlbumNotFound
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
	r, _ := NewPhotosLibraryClient(authenticatedClient)
	return r
}

func defaultCache() Cache {
	return NewCachitaCache()
}

// NewCachedAlbumsService returns an albums Google Photos client with cached results.
// The authenticatedClient should have all oAuth credentials in place.
func NewCachedAlbumsService(authenticatedClient *http.Client, options ...Option) *CachedAlbumsService {
	var repo Repository = defaultRepo(authenticatedClient)
	var albumCache Cache = defaultCache()

	for _, o := range options {
		switch o.Name() {
		case optKeyRepo:
			repo = o.Value().(Repository)
		case optKeyCache:
			albumCache = o.Value().(Cache)
		}
	}

	return &CachedAlbumsService{
		repo:  repo,
		cache: albumCache,
	}
}

const (
	optKeyRepo  = "repository"
	optKeyCache = "cache"
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
func WithRepository(s Repository) option {
	return option{
		name:  optKeyRepo,
		value: s,
	}
}

// WithCache configures the cache.
func WithCache(s Cache) option {
	return option{
		name:  optKeyCache,
		value: s,
	}
}
