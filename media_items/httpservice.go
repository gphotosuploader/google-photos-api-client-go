package media_items

import (
	"context"
	"net/http"

	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/media_items"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
)

// SharedAlbumsService represents a Google Photos client for media management.
type MediaItemsService interface {
	CreateItemsFromFiles(ctx context.Context, albumId string, paths []string, position albums.AlbumPosition) ([]NewMediaItemResult, error)
}

type NewMediaItemResult = media_items.NewMediaItemResult

// repository represents the Google Photos API client for media management.
type repository interface {
	BatchCreateItemsFromFiles(albumId string, paths []string, position albums.AlbumPosition, ctx context.Context) ([]NewMediaItemResult, error)
}

func defaultRepo(authenticatedClient *http.Client, upldr uploader.MediaUploader) media_items.HttpMediaItemsService {
	return media_items.NewHttpMediaItemsService(authenticatedClient, upldr)
}

// HttpMediaItemsService implements a Google Photos client.
type HttpMediaItemsService struct {
	repo repository
}

// CreateItemsFromFiles create one or multiple media items after uploading the file.
func (s HttpMediaItemsService) CreateItemsFromFiles(ctx context.Context, albumId string, paths []string, position albums.AlbumPosition) ([]NewMediaItemResult, error) {
	return s.repo.BatchCreateItemsFromFiles(albumId, paths, position, ctx)
}

