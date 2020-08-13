package gphotos_test

import (
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
)

func TestNewOAuthConfig(t *testing.T) {
	creds := gphotos.APIAppCredentials{
		ClientID:     "test-client-ID",
		ClientSecret: "test-client-secret",
	}

	got := gphotos.NewOAuthConfig(creds)

	if got.ClientID != creds.ClientID {
		t.Errorf("client ID should be equal: got=%s, want=%s", got.ClientID, creds.ClientID)
	}

	if got.ClientSecret != creds.ClientSecret {
		t.Errorf("client secret should be equal: got=%s, want=%s", got.ClientSecret, creds.ClientSecret)
	}
}
