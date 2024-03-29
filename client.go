package gphotos

import (
	"errors"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"net/http"
)

const (
	Version = "v3.0.0"

	defaultBaseURL   = "https://photoslibrary.googleapis.com/"
	defaultUserAgent = "gphotos" + "/" + Version
)

// A Client manages communication with the Google Photos API.
type Client struct {
	// Uploader implementation used when uploading files to Google Photos.
	Uploader MediaUploader

	// Services used for talking to different parts of the Google Photos API.
	Albums     AlbumsService
	MediaItems MediaItemsService
}

// NewClient returns a new Google Photos API client.
// API methods require authentication, provide an [net/http.Client]
// that will perform the authentication for you (such as that provided
// by the [golang.org/x/oauth2] library).
func NewClient(httpClient *http.Client) (*Client, error) {
	return NewClientWithBaseURL(httpClient, defaultBaseURL)
}

// NewClientWithBaseURL returns a new Google Photos API client with a custom baseURL.
// See [NewClient] for more details.
func NewClientWithBaseURL(httpClient *http.Client, baseURL string) (*Client, error) {
	if httpClient == nil {
		return nil, errors.New("client is nil")
	}

	if baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}

	httpClient = addRetryHandler(httpClient)

	// Create the Albums Service using default values.
	albumsConfig := albums.Config{
		Client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	albumsService, err := albums.New(albumsConfig)
	if err != nil {
		return nil, err
	}

	// Create the Media Items Service using default values.
	mediaItemsConfig := media_items.Config{
		Client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
	}
	mediaItemsService, err := media_items.New(mediaItemsConfig)
	if err != nil {
		return nil, err
	}

	simpleUploader, err := uploader.NewSimpleUploader(httpClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		Uploader:   simpleUploader,
		Albums:     albumsService,
		MediaItems: mediaItemsService,
	}, nil
}
