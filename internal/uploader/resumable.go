package uploader

import (
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
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
	logger := defaultLogger()
	storer := defaultStorer()
	endpoint := defaultEndpoint()

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optKeySessionStorer:
			storer = o.Value().(UploadSessionStore)
		case optKeyEndpoint:
			endpoint = o.Value().(string)
		}
	}

	u := &Uploader{
		client: client,
		url:    endpoint,
		resume: false,
		store:  storer,
		log:    logger,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}


// Validate validates the configuration of the Client.
func (u *Uploader) Validate() error {
	if u.resume && u.store == nil {
		return ErrNilStore
	}

	return nil
}

