package gphotos

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/basic"
)

// Config holds the configuration parameters for the client.
type Config struct {
	// Client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
	Client *http.Client

	// Uploader implements file uploads to Google Photos.
	// This package offers two basic.NewBasicUploader() or resumable.NewResumableUploader().
	Uploader MediaUploader

	// AlbumManager implements the Google Photos' album manager.
	AlbumManager AlbumsService

	// MediaItemManager implements the Google Photos' media item manager.
	MediaItemManager MediaItemsService
}

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Uploader   MediaUploader
	Albums     AlbumsService
	MediaItems MediaItemsService
}

// NewClient creates a new instance of Client with the provided configuration.
func NewClient(config Config) (*Client, error) {
	c := &Client{
		Uploader:   config.Uploader,
		Albums:     config.AlbumManager,
		MediaItems: config.MediaItemManager,
	}

	if config.Client == nil {
		return nil, fmt.Errorf("Error")
	}

	client := newRetryHandler(config.Client)

	if c.Albums == nil {
		c.Albums, _ = albums.NewService(albums.Config{Client: client})
	}

	if c.MediaItems == nil {
		c.MediaItems, _ = media_items.NewHttpMediaItemsService(client)
	}

	if c.Uploader == nil {
		c.Uploader, _ = basic.NewBasicUploader(client)
	}

	return c, nil
}

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
