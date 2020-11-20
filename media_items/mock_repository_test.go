package media_items

import "context"

// MockedRepository mocks the repository.
type MockedRepository struct {
	CreateManyFn        func(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	CreateManyToAlbumFn func(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	GetFn               func(ctx context.Context, itemId string) (*MediaItem, error)
	ListByAlbumFn       func(ctx context.Context, albumId string) ([]MediaItem, error)
}

// CreateMany invokes the mock implementation.
func (r MockedRepository) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyFn(ctx, mediaItems)
}

// CreateManyToAlbum invokes the mock implementation.
func (r MockedRepository) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return r.CreateManyToAlbumFn(ctx, albumId, mediaItems)
}

// Get invokes the mock implementation.
func (r MockedRepository) Get(ctx context.Context, itemId string) (*MediaItem, error) {
	return r.GetFn(ctx, itemId)
}

// ListByAlbum invokes the mock implementation.
func (r MockedRepository) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	return r.ListByAlbumFn(ctx, albumId)
}
