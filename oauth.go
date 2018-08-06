package gphotosclient

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

func NewOAuthConfig(creds APIAppCredentials) oauth2.Config {
	return oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint:     google.Endpoint,
	}
}

// NewOAuthClient creates a new http.Client with a bearer access token
// func NewOAuthClient() (*oauth2ns.AuthorizedClient, error) {
// 	photosClient := oauth2ns.Authorize(&OAuthConfig)
// 	return photosClient, nil
// }

// func NewOauthClientFromToken(token *oauth2.Token) *http.Client {
// 	return OAuthConfig.Client(context.Background(), token)
// }
