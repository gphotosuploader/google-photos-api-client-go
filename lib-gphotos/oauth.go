package gphotos

import (
	photoslibrary "github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type APIAppCredentials struct {
	ClientID     string
	ClientSecret string
}

func NewOAuthConfig(creds APIAppCredentials) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint:     google.Endpoint,
	}
}
