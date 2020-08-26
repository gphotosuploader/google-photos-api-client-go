package gphotos_test

import (
	"net/http"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/gphotos"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewGooglePhotosService(t *testing.T) {
	c := http.DefaultClient

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := gphotos.NewGooglePhotosService(c)
		if err != nil {
			t.Errorf("no error expected at this point. err: %s", err)
		}
	})

	t.Run("WithLoggerOption", func(t *testing.T) {
		l := &log.DiscardLogger{}
		_, err := gphotos.NewGooglePhotosService(c, gphotos.WithLogger(l))
		if err != nil {
			t.Errorf("no error expected at this point. err: %s", err)
		}
	})
}