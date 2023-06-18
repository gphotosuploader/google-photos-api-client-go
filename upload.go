package gphotos

import (
	"context"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
)

// Upload uploads the specified file and creates the media item
// in Google Photos.
func (c *Client) Upload(ctx context.Context, filePath string) (*media_items.MediaItem, error) {
	token, err := c.Uploader.UploadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	return c.MediaItems.Create(ctx, media_items.SimpleMediaItem{
		UploadToken: string(token),
	})
}

// UploadToAlbum uploads the specified file and creates the media item
// in the specified album in Google Photos.
func (c *Client) UploadToAlbum(ctx context.Context, albumId string, filePath string) (*media_items.MediaItem, error) {
	token, err := c.Uploader.UploadFile(ctx, filePath)
	if err != nil {
		return nil, err
	}
	item := media_items.SimpleMediaItem{
		UploadToken: string(token),
	}
	return c.MediaItems.CreateToAlbum(ctx, albumId, item)
}
