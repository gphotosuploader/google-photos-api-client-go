package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mock"
)

func TestNewClient(t *testing.T) {
	c := http.DefaultClient

	t.Run("EmptyHTTPClient", func(t *testing.T) {
		_, err := gphotos.NewClient(nil)
		if err == nil {
			t.Errorf("error was expected here")
		}
	})

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := gphotos.NewClient(c)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		l := &log.DiscardLogger{}
		if _, err := gphotos.NewClient(c, gphotos.WithLogger(l)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithCacher", func(t *testing.T) {
		cacher := &mock.Cache{}
		if _, err := gphotos.NewClient(c, gphotos.WithCacher(cacher)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithSessionStorer", func(t *testing.T) {
		s := &mock.SessionStorer{}
		if _, err := gphotos.NewClient(c, gphotos.WithSessionStorer(s)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithPhotoService", func(t *testing.T) {
		s := &mock.PhotoService{}
		if _, err := gphotos.NewClient(c, gphotos.WithPhotoService(s)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("WithUploader", func(t *testing.T) {
		u := &mock.Uploader{}
		if _, err := gphotos.NewClient(c, gphotos.WithUploader(u)); err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}
