package media_items

import (
	"context"
	"errors"
	"net/http"
)

// Repository represents a media items repository.
type Repository interface {
	CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error)
	Get(ctx context.Context, itemId string) (*MediaItem, error)
	ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error)
}

// HttpMediaItemsService implements a media items Google Photos client.
type HttpMediaItemsService struct {
	repo Repository
}

var (
	// NullMediaItem is a zero value MediaItem.
	NullMediaItem = MediaItem{}

	ErrNotFound     = errors.New("media item not found")
	ErrServerFailed = errors.New("internal server error")
)

// Create creates a media item in the repository.
// By default, the media item will be added to the end of the library.
func (ms HttpMediaItemsService) Create(ctx context.Context, mediaItem SimpleMediaItem) (MediaItem, error) {
	result, err := ms.CreateMany(ctx, []SimpleMediaItem{mediaItem})
	if err != nil {
		return NullMediaItem, err
	}
	return result[0], nil
}

// CreateMany creates one or more media items in the repository.
// By default, the media item(s) will be added to the end of the library.
func (ms HttpMediaItemsService) CreateMany(ctx context.Context, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return ms.repo.CreateMany(ctx, mediaItems)
}

// CreateToAlbum creates a media item in the repository.
// If an album id is specified, the media item is also added to the album.
// By default, the media item will be added to the end of the library or album.
func (ms HttpMediaItemsService) CreateToAlbum(ctx context.Context, albumId string, mediaItem SimpleMediaItem) (MediaItem, error) {
	result, err := ms.CreateManyToAlbum(ctx, albumId, []SimpleMediaItem{mediaItem})
	if err != nil {
		return NullMediaItem, err
	}
	return result[0], nil
}

// CreateManyToAlbum creates one or more media item(s) in the repository.
// If an album id is specified, the media item(s) are also added to the album.
// By default, the media item(s) will be added to the end of the library or album.
func (ms HttpMediaItemsService) CreateManyToAlbum(ctx context.Context, albumId string, mediaItems []SimpleMediaItem) ([]MediaItem, error) {
	return ms.repo.CreateManyToAlbum(ctx, albumId, mediaItems)
}

// Get returns the media item specified based on a given media item id.
// It will return ErrNotFound if the media item id does not exist.
func (ms HttpMediaItemsService) Get(ctx context.Context, mediaItemId string) (*MediaItem, error) {
	mediaItem, err := ms.repo.Get(ctx, mediaItemId)
	if err != nil {
		return &NullMediaItem, ErrNotFound
	}
	return mediaItem, nil
}

// ListByAlbum list all media items in the specified album.
func (ms HttpMediaItemsService) ListByAlbum(ctx context.Context, albumId string) ([]MediaItem, error) {
	return ms.repo.ListByAlbum(ctx, albumId)
}

// NewHttpMediaItemsService returns a media items Google Photos client.
// The authenticatedClient should have all oAuth credentials in place.
func NewHttpMediaItemsService(authenticatedClient *http.Client) (*HttpMediaItemsService, error) {
	c, err := NewPhotosLibraryClient(authenticatedClient)
	if err != nil {
		return nil, err
	}
	return &HttpMediaItemsService{
		repo: c,
	}, nil
}
