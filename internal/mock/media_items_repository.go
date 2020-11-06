package mock

import (
	"context"

	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/media_items"
)

// MediaItemsRepository mocks the a media_items repository.
type MediaItemsRepository struct {
	BatchCreateItemsFromFilesFn      func(albumId string, paths []string, position albums.AlbumPosition, ctx context.Context) ([]media_items.NewMediaItemResult, error)
	BatchCreateItemsFromFilesInvoked bool
}

// Get invokes the mock implementation and marks the function as invoked.
func (s MediaItemsRepository) BatchCreateItemsFromFiles(albumId string, paths []string, position albums.AlbumPosition, ctx context.Context) ([]media_items.NewMediaItemResult, error) {
	s.BatchCreateItemsFromFilesInvoked = true
	return s.BatchCreateItemsFromFilesFn(albumId, paths, position, ctx)
}
