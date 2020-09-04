package gphotos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

// UploadItem represents an uploadable item.
type UploadItem interface {
	uploader.UploadItem
}

// FileUploadItem represents a local file.
type FileUploadItem string

// Open returns a stream.
// Caller should close it finally.
func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	f, err := os.Stat(m.String())
	if err != nil {
		return nil, 0, err
	}
	r, err := os.Open(m.String())
	if err != nil {
		return nil, 0, err
	}
	return r, f.Size(), nil
}

// Name returns the filename.
func (m FileUploadItem) Name() string {
	return path.Base(m.String())
}

func (m FileUploadItem) String() string {
	return string(m)
}

func (m FileUploadItem) Size() int64 {
	f, err := os.Stat(m.String())
	if err != nil {
		return 0
	}
	return f.Size()
}

// AddMediaToAlbum returns MediaItem created after uploading the item to Google Photos library.
func (c *Client) AddMediaToLibrary(ctx context.Context, item UploadItem) (*photoslibrary.MediaItem, error) {
	return c.AddMediaToAlbum(ctx, item, "")
}

// AddMediaToAlbum returns MediaItem created after uploading the item and adding it to an`albumID`.
func (c *Client) AddMediaToAlbum(ctx context.Context, item UploadItem, albumID string) (*photoslibrary.MediaItem, error) {
	c.log.Debugf("Initiating upload and media item creation: file=%s", item.String())

	token, err := c.uploader.Upload(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed getting upload token for %s: err=%s", item.String(), err)
	}

	c.log.Debugf("File has been uploaded: file=%s", item.String())

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
		c.log.Errorf("Failed to create media item: file=%s, err=%s", item.String(), err)
		return nil, fmt.Errorf("error while trying to create this media item, err=%s", err)
	}

	return firstMediaItemResult(res.NewMediaItemResults)
}

func firstMediaItemResult(res []*photoslibrary.NewMediaItemResult) (*photoslibrary.MediaItem, error) {
	if len(res) == 0 {
		return nil, nil
	}

	r := res[0]
	if r.Status == nil {
		return nil, errors.New("found unknown error on MediaItem")
	}

	// Google Photos API uses a GRPC code. Values can be obtained at
	// https://godoc.org/google.golang.org/genproto/googleapis/rpc/code
	if r.Status.Code == 0 {
		return r.MediaItem, nil
	}

	return nil, fmt.Errorf("found error on MediaItem: err=%v", r.Status.Message)
}
