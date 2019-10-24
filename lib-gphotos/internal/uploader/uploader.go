package uploader

import (
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos/internal/log"
)

const (
	// API endpoint URL for upload media
	uploadEndpoint = "https://photoslibrary.googleapis.com/v1/uploads"
)

// Uploader is a client for uploading media to Google Photos.
// Original photos library does not provide `/v1/uploads` API.
type Uploader struct {
	// HTTP Client
	client *http.Client
	// URL of the endpoint to upload to
	url string
	// If Resume is true the UploadSessionStore is required.
	resume bool
	// store keeps upload session information.
	store UploadSessionStore

	log log.Logger
}

// NewUploader returns an Uploader using the specified client or error in case
// of non valid configuration.
// The client must have the proper permissions to upload files.
//
// Use WithResumableUploads(...), WithLogger(...) and WithEndpoURL(...) to
// customize configuration.
func NewUploader(client *http.Client, options ...Option) (*Uploader, error) {
	l := log.DiscardLogger{}
	u := &Uploader{
		client: client,
		url:    uploadEndpoint,
		resume: false,
		store:  nil,
		log:    &l,
	}

	for _, opt := range options {
		opt(u)
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil

}

// WithResumableUploads enables resumable uploads.
// Resumable uploads needs an UploadSessionStore to keep upload session information.
func WithResumableUploads(store UploadSessionStore) Option {
	return func(c *Uploader) {
		c.resume = true
		c.store = store
	}
}

// WithLogger sets the logger to log messages.
func WithLogger(l log.Logger) Option {
	return func(c *Uploader) {
		c.log = l
	}
}

// WithEndpointURL sets the URL of the endpoint to upload to.
func WithEndpointURL(url string) Option {
	return func(c *Uploader) {
		c.url = url
	}
}

// Validate validates the configuration of the Client.
func (u *Uploader) Validate() error {
	if u.resume && u.store == nil {
		return ErrNilStore
	}

	return nil
}

// Option defines an option for a Client
type Option func(*Uploader)
