package gphotos

import (
	"fmt"
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader/basic"
)

// Config holds the configuration parameters for the client.
type Config struct {
	// Client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
	Client *http.Client

	// Uploader implements file uploads to Google Photos.
	// This package offers two basic.NewBasicUploader() or resumable.NewResumableUploader().
	Uploader Uploader

	// AlbumService implements the Google Photos' album service.
	AlbumService AlbumsService

	// MediaItemService implements the Google Photos' media item service.
	MediaItemService MediaItemsService
}

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Uploader   Uploader
	Albums     AlbumsService
	MediaItems MediaItemsService
}

// New creates a new instance of Client with the provided configuration.
func New(config Config) (*Client, error) {
	c := &Client{
		Uploader:   config.Uploader,
		Albums:     config.AlbumService,
		MediaItems: config.MediaItemService,
	}

	// config.Client is required unless other services are submitted.
	if config.Client == nil &&
		(config.AlbumService == nil || config.MediaItemService == nil || config.Uploader == nil) {
		return nil, fmt.Errorf("an HTTP client is necessary")
	}

	client := newRetryHandler(config.Client)

	if c.Albums == nil {
		c.Albums, _ = albums.New(albums.Config{Client: client})
	}

	if c.MediaItems == nil {
		c.MediaItems, _ = media_items.New(media_items.Config{Client: client})
	}

	if c.Uploader == nil {
		c.Uploader, _ = basic.NewBasicUploader(client)
	}

	return c, nil
}
