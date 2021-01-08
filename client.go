package gphotos

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

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

// clientWithRetryPolicy returns a HTTP client with a retry policy.
func clientWithRetryPolicy(authenticatedClient *http.Client) *http.Client {
	client := retryablehttp.NewClient()
	client.Logger = nil // Disable DEBUG logs
	client.HTTPClient = authenticatedClient
	return client.StandardClient()
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
// Use WithUploader(), WithAlbumsService(), WithMediaItemsService() to customize it.
//
// There is a resumable uploader implemented on uploader.NewResumableUploader().
func NewClient(authenticatedClient *http.Client, options ...Option) (*Client, error) {
	client, err := defaultGPhotosClient(authenticatedClient)
	if err != nil {
		return nil, err
	}

	for _, o := range options {
		switch o.Name() {
		case optkeyUploader:
			client.Uploader = o.Value().(MediaUploader)
		case optkeyAlbumsService:
			client.Albums = o.Value().(AlbumsService)
		case optkeyMediaItemsService:
			client.MediaItems = o.Value().(MediaItemsService)
		}
	}

	return client, nil
}

const (
	optkeyUploader          = "uploader"
	optkeyAlbumsService     = "albumService"
	optkeyMediaItemsService = "mediaItemsService"
)

// Option represents a configurable parameter.
type Option interface {
	Name() string
	Value() interface{}
}

type option struct {
	name  string
	value interface{}
}

func (o option) Name() string       { return o.name }
func (o option) Value() interface{} { return o.value }

// WithUploader configures the Media Uploader.
func WithUploader(s MediaUploader) *option {
	return &option{
		name:  optkeyUploader,
		value: s,
	}
}

// WithAlbumsService configures the Albums Service.
func WithAlbumsService(s AlbumsService) *option {
	return &option{
		name:  optkeyAlbumsService,
		value: s,
	}
}

// WithMediaItemsService configures the Media Items Service.
func WithMediaItemsService(s MediaItemsService) *option {
	return &option{
		name:  optkeyMediaItemsService,
		value: s,
	}
}
