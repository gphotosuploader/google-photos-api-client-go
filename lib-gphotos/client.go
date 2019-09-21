package gphotos

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// Client is a client for uploading a media.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Client struct {
	*photoslibrary.Service
	*http.Client
	token *oauth2.Token

	log *log.Logger
}

// Token returns the value of the token used by the gphotos Client
// Cannot be used to set the token
func (c *Client) Token() *oauth2.Token {
	if c.token == nil {
		return nil
	}
	return &(*c.token)
}

// NewClient constructs a new PhotosClient from an oauth httpClient
func NewClient(oauthHTTPClient *http.Client, maybeToken ...*oauth2.Token) (*Client, error) {
	var token *oauth2.Token

	if len(maybeToken) > 1 {
		token = maybeToken[0]
	}

	photosService, err := photoslibrary.New(oauthHTTPClient)
	if err != nil {
		return nil, err
	}

	c := Client{
		Service: photosService,
		Client:  oauthHTTPClient,
		token:   token,
		log:     log.New(os.Stdout, logPrefix, logFlags),
	}
	return &c, nil
}
