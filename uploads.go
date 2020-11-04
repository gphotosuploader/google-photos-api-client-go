package gphotos
/*
import (
	"context"
	"errors"
	"fmt"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

// UploadItem represents an uploadable item.
type UploadItem interface {
	uploader.UploadItem
}

// AddMediaToLibrary returns MediaItem created after uploading the item to Google Photos library.
func (c *Client) AddMediaToLibrary(ctx context.Context, item UploadItem) (*photoslibrary.MediaItem, error) {
	return c.AddMediaToAlbum(ctx, item, nil)
}

// AddMediaToAlbum returns MediaItem created after uploading the item and adding it to an`albumID`.
func (c *Client) AddMediaToAlbum(ctx context.Context, item UploadItem, album *photoslibrary.Album) (*photoslibrary.MediaItem, error) {
	c.log.Debugf("Initiating upload and media item creation: file=%s", item.Name())

	var albumID string
	if album != nil {
		albumID = album.Id
	}

	token, err := c.uploader.Upload(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed getting upload token for %s: err=%w", item.Name(), err)
	}

	c.log.Debugf("File has been uploaded: file=%s", item.Name())

	res, err := c.service.CreateMediaItems(ctx, &photoslibrary.BatchCreateMediaItemsRequest{
		AlbumId: albumID,
		NewMediaItems: []*photoslibrary.NewMediaItem{
			{
				Description:     item.Name(),
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: string(token)},
			},
		},
	})
	if err != nil {
		c.log.Errorf("Failed to create media item: file=%s, err=%s", item.Name(), err)
		return nil, fmt.Errorf("error while trying to create this media item, err=%w", err)
	}

	return firstMediaItemResult(res.NewMediaItemResults)
}

func firstMediaItemResult(res []*photoslibrary.NewMediaItemResult) (*photoslibrary.MediaItem, error) {
	if len(res) == 0 || res[0].Status == nil {
		return nil, errors.New("found error on MediaItem")
	}

	// Google Photos API uses a GRPC code. Values can be obtained at
	// https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/code
	if res[0].Status.Code == 0 {
		return res[0].MediaItem, nil
	}

	return nil, fmt.Errorf("found error on MediaItem. err: %s", res[0].Status.Message)
}
*/