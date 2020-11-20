package media_items

import (
	"context"
)

// MediaItemsService represents a Google Photos client for media management.
type MediaItemsService interface {
	Create(ctx context.Context, mediaItem SimpleMediaItem) (MediaItem, error)
	CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	CreateToAlbum(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (MediaItem, error)
	CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	Get(ctx context.Context, mediaItemId string) (*MediaItem, error)
	ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error)
}

// Repository represents a media items repository.
type Repository interface {
	CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	Get(ctx context.Context, itemId string) (*MediaItem, error)
	ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error)
}

// HttpMediaItemsService implements a Google Photos client.
type HttpMediaItemsService struct {
	repo Repository
}

func (ms HttpMediaItemsService) Create(ctx context.Context, mediaItem SimpleMediaItem) (MediaItem, error) {
	result, err := ms.CreateMany(ctx, []SimpleMediaItem{mediaItem})
	if err != nil {
		return MediaItem{}, err
	}
	return result[0], nil
}

func (ms HttpMediaItemsService) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return ms.repo.CreateMany(ctx, mediaItems)
}

func (ms HttpMediaItemsService) CreateToAlbum(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (MediaItem, error) {
	result, err := ms.CreateManyToAlbum(ctx, albumId, []SimpleMediaItem{mediaItem})
	if err != nil {
		return MediaItem{}, err
	}
	return result[0], nil
}

func (ms HttpMediaItemsService) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return ms.repo.CreateManyToAlbum(ctx, albumId, mediaItems)
}

func (ms HttpMediaItemsService) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	return ms.Get(ctx, mediaItemId)
}

func (ms HttpMediaItemsService) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	return ms.ListByAlbum(ctx, albumId)
}
