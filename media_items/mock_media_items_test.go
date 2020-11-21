package media_items

import "context"

type MockedMediaItemsService struct {
	CreateFn func(ctx context.Context, mediaItem SimpleMediaItem) (MediaItem, error)
	CreateManyFn func(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	CreateToAlbumFn func(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (MediaItem, error)
	CreateManyToAlbumFn func(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	GetFn func(ctx context.Context, mediaItemId string) (*MediaItem, error)
	ListByAlbumFn func(ctx context.Context, albumId string) ([]MediaItem, error)
}
