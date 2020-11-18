package media_items

import (
	"context"
	"net/http"

	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/duffpl/google-photos-api-client/media_items"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
)

// MediaItemsService represents a Google Photos client for media management.
type MediaItemsService interface {
	CreateFromFile(ctx context.Context, albumId string, path string) (MediaItem, error)
	CreateManyFromFiles(ctx context.Context, albumId string, paths []string) ([]MediaItem, error)
	Get(ctx context.Context, itemId string) (*MediaItem, error)
	List(ctx context.Context, options *ListOptions, pageToken string) (mediaItems []MediaItem, nextPageToken string, err error)
	Update(ctx context.Context, mediaItem MediaItem, updateMask []Field) (*MediaItem, error)
}

// Repository represents the Google Photos API client for media management.
type Repository interface {
	BatchCreateItemsFromFiles(albumId string, paths []string, position albums.AlbumPosition, ctx context.Context) ([]media_items.NewMediaItemResult, error)
}

func defaultRepo(authenticatedClient *http.Client, upldr uploader.MediaUploader) media_items.HttpMediaItemsService {
	return media_items.NewHttpMediaItemsService(authenticatedClient, upldr)
}

// HttpMediaItemsService implements a Google Photos client.
type HttpMediaItemsService struct {
	repo Repository
}

// CreateFromFile create one media items after uploading the file.
// Media item will be added to the end of the Album.
func (s HttpMediaItemsService) CreateFromFile(ctx context.Context, albumId string, path string) (MediaItem, error) {
	nullMediaItem := MediaItem{}
	results, err := s.CreateManyFromFiles(ctx, albumId, []string{path})
	if err != nil {
		return nullMediaItem, err
	}
	return results[0], nil
}

// CreateManyFromFiles create one or multiple media items after uploading the files.
// Media items will be added to the end of the Album.
func (s HttpMediaItemsService) CreateManyFromFiles(ctx context.Context, albumId string, paths []string) ([]MediaItem, error) {
	mediaItems := make([]MediaItem, 0)
	results, err := s.repo.BatchCreateItemsFromFiles(albumId, paths, albums.AlbumPosition{Position: albums.AlbumPositionTypeLastInAlbum}, ctx)
	if err != nil {
		return mediaItems, err
	}
	for _, i := range results {
		mediaItems = append(mediaItems, i.MediaItem)
	}
	return mediaItems, nil
}
