package gphotos

import (
	"net/http"
	"sync"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos/internal/uploader"
)

// Client is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	// Google Photos client
	*photoslibrary.Service
	// Uploader to upload new files to Google Photos
	uploader *uploader.Uploader

	log log.Logger
	mu  sync.Mutex

	token *oauth2.Token // DEPRECATED: `token` will disappear in the next MAJOR version.
}

// NewClientWithResumableUploads constructs a new gphotos.Client from the provided HTTP client and
// the given options.
//
// `httpClient` is an client with authentication credentials.
// `store` is an UploadSessionStore to keep upload sessions to resume uploads.
func NewClientWithResumableUploads(httpClient *http.Client, store uploader.UploadSessionStore, options ...Option) (*Client, error) {
	photosService, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}

	upldr, err := uploader.NewUploader(httpClient, uploader.WithResumableUploads(store))
	if err != nil {
		return nil, err
	}

	c := &Client{
		Service:  photosService,
		uploader: upldr,
		log:      log.NewDiscardLogger(),
	}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

// WithLogger set a new Logger to log messages.
func WithLogger(l log.Logger) func(*Client) {
	return func(c *Client) {
		c.log = l
	}
}

// Option defines an option for a Client
type Option func(*Client)

// codebeat:disable

// NewClient constructs a new PhotosClient from an oauth httpClient.
//
// `httpClient` is an HTTP Client with authentication credentials.
//
// DEPRECATED: Use NewClientWithOptions(...) instead.
// This package doesn't need Client.token anymore, used `Client.Client` instead.
func NewClient(httpClient *http.Client, maybeToken ...*oauth2.Token) (*Client, error) {
	var token *oauth2.Token

	if len(maybeToken) > 0 {
		token = maybeToken[0]
	}

	photosService, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}

	upldr, err := uploader.NewUploader(httpClient)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Service:  photosService,
		uploader: upldr,
		log:      log.NewDiscardLogger(),
	}

	c.token = token
	return c, nil
}

// Token returns the value of the token used by the gphotos Client
// Cannot be used to set the token
//
// DEPRECATED: Use the authenticated HTTP Client `Client.Client` instead.
func (c *Client) Token() *oauth2.Token {
	if c.token == nil {
		return nil
	}
	return &(*c.token)
}

// codebeat:enable
