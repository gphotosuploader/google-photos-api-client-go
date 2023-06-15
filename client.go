package gphotos

import (
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

	// config.Client is required unless other services are submitted.
	if config.Client == nil &&
		(config.AlbumManager == nil || config.MediaItemManager == nil || config.Uploader == nil) {
		return nil, fmt.Errorf("an HTTP client is necessary")
	}

	client := newRetryHandler(config.Client)

	if c.Albums == nil {
		c.Albums, _ = albums.NewService(albums.Config{Client: client})
	}

	if c.MediaItems == nil {
		c.MediaItems, _ = media_items.New(client)
	}

	if c.Uploader == nil {
		c.Uploader, _ = basic.NewBasicUploader(client)
	}

	return c, nil
}
