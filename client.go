package gphotos

import (
	"net/http"

	"github.com/duffpl/google-photos-api-client/media_items"
	"github.com/duffpl/google-photos-api-client/uploader"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
)

// Client is a Google Photos client with enhanced capabilities.
type Client struct {
	Albums     albums.AlbumsService
	MediaItems media_items.MediaItemsService
}

// NewClient constructs a new gphotos.Client from the provided HTTP client and the given options.
// The client is an HTTP client used for calling Google Photos. It needs the proper authentication in place.
//
// By default it will use a in memory cache for Albums repository.
//
// Use WithUploader(), WithAlbumsService() to customize it.
//
// There is a resumable uploader implemented on uploader.NewResumableUploader().
func NewClient(authenticatedClient *http.Client, options ...Option) (*Client, error) {
	var upldr uploader.MediaUploader = uploader.NewHttpMediaUploader(authenticatedClient)
	var albumsService albums.AlbumsService = albums.NewCachedAlbumsService(authenticatedClient)

	for _, o := range options {
		switch o.Name() {
		case optkeyUploader:
			upldr = o.Value().(uploader.MediaUploader)
		case optkeyAlbumsService:
			albumsService = o.Value().(albums.AlbumsService)
		}
	}

	return &Client{
		Albums:     albumsService,
		MediaItems: media_items.NewHttpMediaItemsService(authenticatedClient, upldr),
	}, nil

}

const (
	optkeyUploader      = "uploader"
	optkeyAlbumsService = "albumService"
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
func WithUploader(s uploader.MediaUploader) Option {
	return &option{
		name:  optkeyUploader,
		value: s,
	}
}

// WithAlbumsService configures the Albums Service.
func WithAlbumsService(s albums.AlbumsService) Option {
	return &option{
		name:  optkeyAlbumsService,
		value: s,
	}
}
