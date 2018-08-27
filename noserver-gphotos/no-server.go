package gphotos

import (
	"context"
	"net/url"

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

type AuthenticateUserOption func(*AuthenticateUserFuncConfig) error
type AuthenticateUserFuncConfig struct {
	userLoginHint string
}

func WithUserLoginHint(loginHint string) AuthenticateUserOption {
	return func(config *AuthenticateUserFuncConfig) error {
		config.userLoginHint = loginHint
		return nil
	}
}

// AuthenticateUser() option creates a new http.Client with a bearer access token
func AuthenticateUser(oauthConfig *oauth2.Config, options ...AuthenticateUserOption) ClientConstructorOption {
	var funcConfig AuthenticateUserFuncConfig
	for _, optionFunc := range options {
		optionFunc(&funcConfig)
	}

	// apply oauth2ns options
	var oauth2nsOptions []oauth2ns.AuthenticateUserOption
	if funcConfig.userLoginHint != "" {
		var urlValues = url.Values{}
		urlValues.Set("login_hint", funcConfig.userLoginHint)
		oauth2nsOptions = append(oauth2nsOptions, oauth2ns.WithAuthCallHTTPParams(urlValues))
	}

	return func() (*AuthorizedClient, error) {
		authorizedClient, err := oauth2ns.AuthenticateUser(oauthConfig, oauth2nsOptions...)
		if err != nil {
			return nil, stacktrace.Propagate(err, "failed authenticating user")
		}
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
