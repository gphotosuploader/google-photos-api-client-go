package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
)

func TestNewClient(t *testing.T) {
	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := gphotos.NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithUploader", func(t *testing.T) {
		u := &mock.Uploader{}
		if _, err := gphotos.NewClient(http.DefaultClient, gphotos.WithUploader(u)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}
