package gphotos

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/basic"
)

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Albums     AlbumsService
	MediaItems MediaItemsService
	Uploader   uploader.MediaUploader
}

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

// NewClient constructs a new gphotos.Client from the provided HTTP client and the given options.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
//
// By default it will use a in memory cache for Albums repository and implements retries with Exponential backoff.
//
// Use WithUploader(), WithAlbumsService(), WithAlbumsService() to customize it.
//
// There is a resumable uploader implemented on uploader.NewResumableUploader().
func NewClient(authenticatedClient *http.Client, options ...Option) (*Client, error) {
	client := retryablehttp.NewClient()
	client.HTTPClient = authenticatedClient

	var albumsService AlbumsService = albums.NewCachedAlbumsService(client.StandardClient())

	var upldr uploader.MediaUploader
	upldr, err := basic.NewBasicUploader(client.StandardClient())
	if err != nil {
		return nil, err
	}

	var mediaItemsService MediaItemsService
	mediaItemsService, err = media_items.NewHttpMediaItemsService(client.StandardClient())
	if err != nil {
		return nil, err
	}

	for _, o := range options {
		switch o.Name() {
		case optkeyUploader:
			upldr = o.Value().(uploader.MediaUploader)
		case optkeyAlbumsService:
			albumsService = o.Value().(AlbumsService)
		case optkeyMediaItemsService:

		}
	}

	return &Client{
		Albums:     albumsService,
		MediaItems: mediaItemsService,
		Uploader:   upldr,
	}, nil

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
func WithUploader(s uploader.MediaUploader) *option {
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
