package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mocks"
)

func TestNewClient(t *testing.T) {
	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := gphotos.NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}

func TestWithAlbumsService(t *testing.T) {
	want := &mocks.MockedAlbumsService{}

	got, err := gphotos.NewClient(http.DefaultClient, gphotos.WithAlbumsService(want))
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	if got.Albums != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithMediaItemsService(t *testing.T) {
	want := &mocks.MockedMediaItemsService{}

	got, err := gphotos.NewClient(http.DefaultClient, gphotos.WithMediaItemsService(want))
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	if got.MediaItems != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithUploader(t *testing.T) {
	want := &mocks.MockedUploader{}

	got, err := gphotos.NewClient(http.DefaultClient, gphotos.WithUploader(want))
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	if got.Uploader != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
