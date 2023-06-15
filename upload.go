package gphotos

import (
	"context"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
)

// UploadFileToLibrary uploads the specified file to Google Photos.
func (c Client) UploadFileToLibrary(ctx context.Context, filePath string) (media_items.MediaItem, error) {
	token, err := c.Uploader.UploadFile(ctx, filePath)
	if err != nil {
		return media_items.MediaItem{}, err
	}
	return c.MediaItems.Create(ctx, media_items.SimpleMediaItem{
		UploadToken: token,
		FileName:    filePath,
	})
}

// UploadFileToAlbum uploads the specified file to the album in Google Photos.
func (c Client) UploadFileToAlbum(ctx context.Context, albumId string, filePath string) (media_items.MediaItem, error) {
	token, err := c.Uploader.UploadFile(ctx, filePath)
	if err != nil {
		return media_items.MediaItem{}, err
	}
	item := media_items.SimpleMediaItem{
		UploadToken: token,
		FileName:    filePath,
	}
	return c.MediaItems.CreateToAlbum(ctx, albumId, item)
}
