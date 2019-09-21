package gphotos

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/xerrors"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// Client is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	*photoslibrary.Service
	*http.Client

	token *oauth2.Token // DEPRECATED: `token` will disappear in the next MAJOR version.

	log *log.Logger
}

// NewClient constructs a new PhotosClient from an oauth httpClient.
//
// `httpClient` is an HTTP Client with authentication credentials.
//
// `maybeToken` will disappear in the next MAJOR version.
// This package doesn't need Client.token anymore, used `Client.Client` instead.
func NewClient(httpClient *http.Client, maybeToken ...*oauth2.Token) (*Client, error) {
	var token *oauth2.Token

	if httpClient == nil {
		return nil, xerrors.New("client is nil")
	}

	if len(maybeToken) > 0 {
		token = maybeToken[0]
	}

	photosService, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}

	c := Client{
		Service: photosService,
		Client:  httpClient,
		token:   token,
		log:     defaultLogger(),
	}
	return &c, nil
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
