package gphotos

import (
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// APIAppCredentials represents the credentials for Google Photos OAuth client.
type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

// NewOAuthConfig returns the OAuth configuration for Google Photos service.
func NewOAuthConfig(creds APIAppCredentials) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint:     google.Endpoint,
	}
}
