package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mocks"
)

func TestNewClient(t *testing.T) {
	t.Run("Fail without an HTTP client", func(t *testing.T) {
		_, err := gphotos.NewClient(gphotos.Config{})
		if err == nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("Success with only an HTTP client", func(t *testing.T) {
		_, err := gphotos.NewClient(gphotos.Config{Client: http.DefaultClient})
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("Success with a custom AlbumManager", func(t *testing.T) {
		want := &mocks.MockedAlbumsService{}

		c := gphotos.Config{
			Client:       http.DefaultClient,
			AlbumManager: want,
		}
		got, err := gphotos.NewClient(c)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		if got.Albums != want {
			t.Errorf("want: %v, got: %v", want, got.Albums)
		}
	})

	t.Run("Success with a custom MediaItemManager", func(t *testing.T) {
		want := &mocks.MockedMediaItemsService{}

		c := gphotos.Config{
			Client:           http.DefaultClient,
			MediaItemManager: want,
		}
		got, err := gphotos.NewClient(c)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		if got.MediaItems != want {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("Success with a custom Uploader", func(t *testing.T) {
		want := &mocks.MockedUploader{}

		c := gphotos.Config{
			Client:   http.DefaultClient,
			Uploader: want,
		}
		got, err := gphotos.NewClient(c)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		if got.Uploader != want {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("Success with a custom Uploader, AlbumManager, MediaItemManager but without Client", func(t *testing.T) {
		uploader := &mocks.MockedUploader{}
		albumManager := &mocks.MockedAlbumsService{}
		mediaItemManager := &mocks.MockedMediaItemsService{}

		c := gphotos.Config{
			Uploader:         uploader,
			AlbumManager:     albumManager,
			MediaItemManager: mediaItemManager,
		}
		_, err := gphotos.NewClient(c)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}
