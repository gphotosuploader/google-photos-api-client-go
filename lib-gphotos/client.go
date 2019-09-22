package gphotos

import (
	"log"
	"net/http"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/oauth2"
)

// Client is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	// Google Photos client
	*photoslibrary.Service
	// HTTP Client
	*http.Client
	// If Resume is true the Store is required.
	resume bool
	// Store map an upload's fingerprint with the corresponding upload URL.
	// Resume enables resumable upload.
	store Store

	log *log.Logger

	token *oauth2.Token // DEPRECATED: `token` will disappear in the next MAJOR version.
}

// NewClient constructs a new gphotos.Client from the provided HTTP client and
// the given options.
//
// `httpClient` is an client with authentication credentials.
func NewClientWithOptions(httpClient *http.Client, options ...Option) (*Client, error) {
	photosService, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Service: photosService,
		Client:  httpClient,
		resume:  false,
		store:   nil,
		log:     defaultLogger(),
	}

	for _, opt := range options {
		opt(c)
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// Option defines an option for a Client
type Option func(*Client)

// OptionResumeUploads enable resumable uploads for the client.
// Resumable uploads needs a gphotos.Store to keep uploads sessions.
func OptionResumeUploads(store Store) func(*Client) {
	return func(c *Client) {
		c.resume = true
		c.store = store
	}
}

// OptionLog set logging for client.
func OptionLog(l *log.Logger) func(*Client) {
	return func(c *Client) {
		c.log = l
	}
}

// Validate validates the configuration of the Client.
func (c *Client) Validate() error {
	if c.resume && c.store == nil {
		return ErrNilStore
	}

	return nil
}

// CanResume returns if the Client can use resumable uploads
func (c *Client) CanResumeUploads() bool {
	return c.resume
}

// NewClient constructs a new PhotosClient from an oauth httpClient.
//
// `httpClient` is an HTTP Client with authentication credentials.
//
// DEPRECATED: This method will disappear in the next MAJOR version. Use NewClientWithOptions()
// This package doesn't need Client.token anymore, used `Client.Client` instead.
func NewClient(httpClient *http.Client, maybeToken ...*oauth2.Token) (*Client, error) {
	var token *oauth2.Token

	if len(maybeToken) > 0 {
		token = maybeToken[0]
	}

	c, err := NewClientWithOptions(httpClient)
	if err != nil {
		return nil, err
	}

	c.token = token
	return c, nil
}

// Token returns the value of the token used by the gphotos Client
// Cannot be used to set the token
//
// DEPRECATED: This method will disappear in the next MAJOR version.
// Use the authenticated HTTP Client `Client.Client` instead.
func (c *Client) Token() *oauth2.Token {
	if c.token == nil {
		return nil
	}
	return &(*c.token)
}
