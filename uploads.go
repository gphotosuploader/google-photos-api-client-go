package gphotos

import (
	"context"
	"errors"
	"fmt"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// AddMediaItem returns MediaItem created after uploading `filename` and adding it to `albumID`.
func (c *Client) AddMediaItem(ctx context.Context, filename, albumID string) (*photoslibrary.MediaItem, error) {
	c.log.Debugf("Initiating upload and media item creation: file=%s", filename)

	uploadToken, err := c.uploader.UploadFromFile(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("failed getting uploadToken for %s: err=%s", filename, err)
	}

	c.log.Debugf("File has been uploaded: file=%s", filename)

	mediaItem, err := c.createMediaItemFromUploadToken(ctx, uploadToken, albumID, filename)
	if err != nil {
		c.log.Errorf("Failed to create media item: file=%s, err=%s", filename, err)
		return nil, fmt.Errorf("error while trying to create this media item, err=%s", err)
	}

	c.log.Debugf("File uploaded and media item created successfully: file=%s", filename)
	return mediaItem, nil
}

func (c *Client) createMediaItemFromUploadToken(ctx context.Context, uploadToken, albumID, filename string) (*photoslibrary.MediaItem, error) {
	req := photoslibrary.BatchCreateMediaItemsRequest{
		AlbumId: albumID,
		NewMediaItems: []*photoslibrary.NewMediaItem{
			{
				Description:     filename,
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
			},
		},
	}

	res, err := c.retryableMediaItemBatchCreateDo(ctx, &req, filename)
	if err != nil {
		return nil, err
	}

	if res == nil || len(res.NewMediaItemResults) != 1 {
		return nil, errors.New("len(batchResults) should be 1")
	}

	result := res.NewMediaItemResults[0]
	// `result.Status.Code` has the GRPC code returned by Google Photos API. Values can be obtained at
	// https://godoc.org/google.golang.org/genproto/googleapis/rpc/code
	if result.Status.Code != 0 {
		return nil, errors.New(result.Status.Message)
	}
	return result.MediaItem, nil
}
