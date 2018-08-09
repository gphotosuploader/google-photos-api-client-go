package gphotoslib

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
// OAuthConfig = oauth2.Config{
// 	ClientID:     config.API_APP_CREDENTIALS.ClientID,
// 	ClientSecret: config.API_APP_CREDENTIALS.ClientSecret,
// 	Scopes:       []string{photoslibrary.PhotoslibraryScope},
// 	Endpoint:     google.Endpoint,
// }
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
