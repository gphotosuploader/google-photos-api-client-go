package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
)

func TestNewClient(t *testing.T) {
	t.Run("Should success with httpClient", func(t *testing.T) {
		_, err := gphotos.NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("Should fail without httpClient", func(t *testing.T) {
		_, err := gphotos.NewClient(nil)
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})
}

func TestNewClientWithBaseURL(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		_, err := gphotos.NewClientWithBaseURL(http.DefaultClient, "https://foo.bar")
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("Should fail without httpClient", func(t *testing.T) {
		_, err := gphotos.NewClientWithBaseURL(nil, "https://foo.bar")
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})

	t.Run("Should fail without baseURL", func(t *testing.T) {
		_, err := gphotos.NewClientWithBaseURL(http.DefaultClient, "")
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})

}
