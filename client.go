package gphotos

import (
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/basic"
)

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Albums     AlbumsService
	MediaItems MediaItemsService
	Uploader   MediaUploader
}

// defaultGPhotosClient returns a gphotos client using the defaults.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
// By default it will use a in memory cache for Albums repository and implements retries with Exponential backoff.
func defaultGPhotosClient(authenticatedClient *http.Client) (*Client, error) {
	client := clientWithRetryPolicy(authenticatedClient)

	var albumsService AlbumsService = albums.NewCachedAlbumsService(client)

	var upldr MediaUploader
	upldr, err := basic.NewBasicUploader(client)
	if err != nil {
		return nil, err
	}

	var mediaItemsService MediaItemsService
	mediaItemsService, err = media_items.NewHttpMediaItemsService(client)
	if err != nil {
		return nil, err
	}

	return &Client{
		Albums:     albumsService,
		MediaItems: mediaItemsService,
		Uploader:   upldr,
	}, nil
}

// NewClient constructs a new gphotos.Client from the provided HTTP client and the given options.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
//
// By default it will use a in memory cache for Albums repository and implements retries with Exponential backoff.
//
// There is a resumable uploader implemented on uploader.NewResumableUploader().
func NewClient(authenticatedClient *http.Client) (*Client, error) {
	client, err := defaultGPhotosClient(authenticatedClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}
