package gphotos

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"
	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	oauth2ns "github.com/nmrshll/oauth2-noserver"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
)

type AuthorizedClient oauth2ns.AuthorizedClient

type ClientConstructorOption func() (*AuthorizedClient, error)

func NewClient(options ...ClientConstructorOption) (client *Client, merr error) {
	if len(options) == 0 {
		return nil, stacktrace.NewError("NewClient expects at least one option")
	}
	for _, option := range options {
		authorizedClient, err := option()
		if err != nil {
			// aggregate any errors
			merr = multierror.Append(merr, err)
			continue
		}
		gphotosClient, err := gphotos.NewClient(authorizedClient.Client, authorizedClient.Token)
		if err != nil {
			// aggregate any errors
			merr = multierror.Append(merr, err)
			continue
		}
		return &Client{*gphotosClient}, nil
	}
	// if all constructor options failed, returned errors returned by each option
	return nil, stacktrace.Propagate(merr, "all options failed with errors:")
}

// func FromOAuthCredentials(apiAppCredentials gphotos.APIAppCredentials) ClientConstructorOption {
// 	return func() (*oauth2ns.AuthorizedClient, error) {
// 		oauthConfig := gphotos.NewOAuthConfig(apiAppCredentials)
// 		photosClient := oauth2ns.Authorize(&oauthConfig)
// 		return photosClient, nil
// 	}
// }

// AuthenticateUser() option creates a new http.Client with a bearer access token
func AuthenticateUser(oauthConfig *oauth2.Config) ClientConstructorOption {
	return func() (*AuthorizedClient, error) {
		authorizedClient := oauth2ns.Authorize(oauthConfig)
		return (*AuthorizedClient)(authorizedClient), nil
	}
}

func FromToken(oauthConfig *oauth2.Config, token *oauth2.Token) ClientConstructorOption {
	return func() (*AuthorizedClient, error) {
		return &AuthorizedClient{
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
