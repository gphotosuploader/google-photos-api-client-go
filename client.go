package gphotos

import (
	"net/http"

	"github.com/duffpl/google-photos-api-client/media_items"
	"github.com/duffpl/google-photos-api-client/shared_albums"
	"github.com/duffpl/google-photos-api-client/uploader"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
)

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Albums       albums.AlbumsService
	MediaItems   media_items.MediaItemsService
	SharedAlbums shared_albums.SharedAlbumsService
}

// NewClient constructs a new gphotos.Client from the provided HTTP client and the given options.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
//
// Use WithAlbumsService(), WithUploader() to customize it.
func NewClient(authenticatedClient *http.Client, options ...Option) (*Client, error) {
	var upldr uploader.MediaUploader = uploader.NewHttpMediaUploader(authenticatedClient)
	var albumsService albums.AlbumsService = albums.NewCachedAlbumsService(authenticatedClient)

	for _, o := range options {
		switch o.Name() {
		case optkeyUploader:
			upldr = o.Value().(uploader.MediaUploader)
		case optKeyAlbumsService:
			albumsService = o.Value().(albums.AlbumsService)
		}
	}

	return &Client{
		Albums:       albumsService,
		MediaItems:   media_items.NewHttpMediaItemsService(authenticatedClient, upldr),
		SharedAlbums: shared_albums.NewHttpSharedAlbumsService(authenticatedClient),
	}, nil

}
