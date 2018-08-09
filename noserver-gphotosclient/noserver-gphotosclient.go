package gphotosclient

import (
	"context"
	"log"

	"github.com/nmrshll/google-photos-api-client-go/gphotoslib"
	oauth2ns "github.com/nmrshll/oauth2-noserver"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
)

type ClientConstructorOption func() (*oauth2ns.AuthorizedClient, error)

func NewClient(options ...ClientConstructorOption) (*gphotoslib.Client, error) {
	if len(options) == 0 {
		return nil, stacktrace.NewError("NewClient expects at least one option")
	}
	for _, option := range options {
		authorizedClient, err := option()
		if err != nil {
			log.Fatal(err)
			continue
		}
		return gphotoslib.NewClient(authorizedClient.Client)
	}
	return nil, stacktrace.NewError("all options returned errors")
}

// func FromOAuthCredentials(apiAppCredentials gphotoslib.APIAppCredentials) ClientConstructorOption {
// 	return func() (*oauth2ns.AuthorizedClient, error) {
// 		oauthConfig := gphotoslib.NewOAuthConfig(apiAppCredentials)
// 		photosClient := oauth2ns.Authorize(&oauthConfig)
// 		return photosClient, nil
// 	}
// }

// type AuthenticateUserOption func() ()

// AuthenticateUser() option creates a new http.Client with a bearer access token
func AuthenticateUser(oauthConfig *oauth2.Config) ClientConstructorOption {
	return func() (*oauth2ns.AuthorizedClient, error) {
		photosClient := oauth2ns.Authorize(oauthConfig)
		return photosClient, nil
	}
}

func FromToken(oauthConfig *oauth2.Config, token *oauth2.Token) ClientConstructorOption {
	return func() (*oauth2ns.AuthorizedClient, error) {
		return &oauth2ns.AuthorizedClient{
			Client: oauthConfig.Client(context.Background(), token),
			Token:  token,
		}, nil
	}
}

// func NewOAuthClient(oauthCreds) (*oauth2ns.AuthorizedClient, error) {
// 	photosClient := oauth2ns.Authorize(&OAuthConfig)
// 	return photosClient, nil
// }

// func NewOauthClientFromToken(token *oauth2.Token) *http.Client {
// 	return OAuthConfig.Client(context.Background(), token)
// }
