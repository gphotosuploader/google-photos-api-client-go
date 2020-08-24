package gphotos

import (
	"net/http"
	"sync"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
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
