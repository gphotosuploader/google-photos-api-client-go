package mocks

import (
	"context"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
)

// MockedMediaItemsService mocks the media items service.
type MockedMediaItemsService struct {
	CreateFn            func(ctx context.Context, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error)
	CreateManyFn        func(ctx context.Context, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error)
	CreateToAlbumFn     func(ctx context.Context, albumId string, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error)
	CreateManyToAlbumFn func(ctx context.Context, albumId string, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error)
	GetFn               func(ctx context.Context, mediaItemId string) (*media_items.MediaItem, error)
	ListByAlbumFn       func(ctx context.Context, albumId string) ([]media_items.MediaItem, error)
}

// Create invokes the mock implementation.
func (m MockedMediaItemsService) Create(ctx context.Context, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error) {
	return m.CreateFn(ctx, mediaItem)
}

// CreateMany invokes the mock implementation.
func (m MockedMediaItemsService) CreateMany(ctx context.Context, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error) {
	return m.CreateManyFn(ctx, mediaItems)
}

// CreateToAlbum invokes the mock implementation.
func (m MockedMediaItemsService) CreateToAlbum(ctx context.Context, albumId string, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error) {
	return m.CreateToAlbumFn(ctx, albumId, mediaItem)
}

// CreateManyToAlbum invokes the mock implementation.
func (m MockedMediaItemsService) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []media_items.SimpleMediaItem) ([]media_items.MediaItem, error) {
	return m.CreateManyToAlbumFn(ctx, albumId, mediaItems)
}

// Get invokes the mock implementation.
func (m MockedMediaItemsService) Get(ctx context.Context, mediaItemId string) (*media_items.MediaItem, error) {
	return m.GetFn(ctx, mediaItemId)
}

// ListByAlbum invokes the mock implementation.
func (m MockedMediaItemsService) ListByAlbum(ctx context.Context, albumId string) ([]media_items.MediaItem, error) {
	return m.ListByAlbumFn(ctx, albumId)
}
