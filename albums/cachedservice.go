package albums

import (
	"context"
	"errors"
	"net/http"

	"github.com/duffpl/google-photos-api-client/albums"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
)

var (
	ErrAlbumNotFound = errors.New("album not found")
)

type AlbumsService interface {
	albums.AlbumsService
	GetByTitle(title string, ctx context.Context) (*albums.Album, error)
}

type CachedAlbumsService struct {
	albums.AlbumsService
	cache cache.Cache
}

// Create new album
func (s CachedAlbumsService) Create(title string, ctx context.Context) (*albums.Album, error) {
	albumPtr, err := s.AlbumsService.Create(title, ctx)
	if err != nil {
		return nil, err
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// Fetch album by id
func (s CachedAlbumsService) Get(id string, ctx context.Context) (*albums.Album, error) {
	albumPtr, err := s.AlbumsService.Get(id, ctx)
	if err != nil {
		return nil, err
	}
	return albumPtr, s.cache.PutAlbum(ctx, *albumPtr)
}

// Fetch album by title. Caches all the albums seen until return the matching one.
func (s CachedAlbumsService) GetByTitle(title string, ctx context.Context) (*albums.Album, error) {
	album, err := s.cache.GetAlbum(ctx, title)
	if err == nil {
		return &album, nil // album was found in the cache
	}

	albumsC, errorsC := s.AlbumsService.ListAllAsync(&albums.AlbumsListOptions{}, ctx)
	for {
		select {
		case item, ok := <-albumsC:
			if !ok {
				return &albums.Album{}, ErrAlbumNotFound // there aren't more albums, album not found
			}
			if item.Title == title {
				return &item, s.cache.PutAlbum(ctx, item) // found, cache it and return it
			}
			if err := s.cache.PutAlbum(ctx, item); err != nil {
				return &albums.Album{}, err // error when caching, returns
			}
		case err := <-errorsC:
			return nil, err
		}
	}
}

// Synchronous wrapper for ListAllAsync with caching. Caches all the returned albums.
func (s CachedAlbumsService) ListAll(options *albums.AlbumsListOptions, ctx context.Context) ([]albums.Album, error) {
	result := make([]albums.Album, 0)
	if err := s.cache.InvalidateAllAlbums(ctx); err != nil {
		return result, err
	}
	albumsC, errorsC := s.AlbumsService.ListAllAsync(options, ctx)
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
			return nil, err
		}
	}
}

// Patches album. updateMask argument can be used to update only selected fields. Currently only id, title
// and coverPhotoMediaItemId are read
func (s CachedAlbumsService) Patch(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error) {
	if err := s.cache.InvalidateAlbum(ctx, album.Title); err != nil {
		return nil, err
	}
	albumPtr, err := s.AlbumsService.Patch(album, updateMask, ctx)
	if err != nil {
		return nil, err
	}
	err = s.cache.PutAlbum(ctx, *albumPtr)
	return albumPtr, err
}

func NewCachedAlbumsService(authenticatedClient *http.Client, options ...Option) CachedAlbumsService {
	var albumsAPIClient albums.AlbumsService = albums.NewHttpAlbumsService(authenticatedClient)
	var albumCache cache.Cache = cache.NewCachitaCache()

	for _, o := range options {
		switch o.Name() {
		case optkeyAlbumsAPIClient:
			albumsAPIClient = o.Value().(albums.AlbumsService)
		case optkeyCacher:
			albumCache = o.Value().(cache.Cache)
		}
	}

	return CachedAlbumsService{
		AlbumsService: albumsAPIClient,
		cache:         albumCache,
	}
}

const (
	optkeyAlbumsAPIClient = "albumsAPIClient"
	optkeyCacher          = "cacher"
)

// Option represents a configurable parameter for Google Photos API client.
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

// WithAlbumsAPIClient configures a Google Photos service
func WithAlbumsAPIClient(s albums.AlbumsService) Option {
	return &option{
		name:  optkeyAlbumsAPIClient,
		value: s,
	}
}

// WithCache configures a cache
func WithCacher(s cache.Cache) Option {
	return &option{
		name:  optkeyCacher,
		value: s,
	}
}
