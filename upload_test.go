package gphotos_test

import (
	"context"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"net/http"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	t.Run("Should success with valid file", func(t *testing.T) {
		srv := mocks.NewMockedGooglePhotosService()
		defer srv.Close()

		httpClient := http.DefaultClient

		mockedUploader, err := uploader.NewSimpleUploader(httpClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		mockedUploader.BaseURL = srv.URL() + "/v1/uploads"

		mediaItemsConfig := media_items.Config{
			Client:  httpClient,
			BaseURL: srv.URL(),
		}
		mockedMediaItems, err := media_items.New(mediaItemsConfig)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		client, err := gphotos.NewClient(httpClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		client.Uploader = mockedUploader
		client.MediaItems = mockedMediaItems

		mediaItem, err := client.Upload(context.Background(), "testdata/upload-success")
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		want := mocks.UploadToken + "Id"
		if want != mediaItem.ID {
			t.Errorf("want: %s, got: %s", want, mediaItem.ID)
		}
	})

	t.Run("Should fail with invalid file", func(t *testing.T) {
		client, err := gphotos.NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		_, err = client.Upload(context.Background(), "non-existent")
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})

}

func TestClient_UploadToAlbum(t *testing.T) {
	t.Run("Should success with valid file", func(t *testing.T) {
		srv := mocks.NewMockedGooglePhotosService()
		defer srv.Close()

		httpClient := http.DefaultClient

		mockedUploader, err := uploader.NewSimpleUploader(httpClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		mockedUploader.BaseURL = srv.URL() + "/v1/uploads"

		albumsConfig := albums.Config{
			Client:  httpClient,
			BaseURL: srv.URL(),
		}
		mockedAlbums, err := albums.New(albumsConfig)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		mediaItemsConfig := media_items.Config{
			Client:  httpClient,
			BaseURL: srv.URL(),
		}
		mockedMediaItems, err := media_items.New(mediaItemsConfig)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		client, err := gphotos.NewClient(httpClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
		client.Uploader = mockedUploader
		client.MediaItems = mockedMediaItems
		client.Albums = mockedAlbums

		mediaItem, err := client.UploadToAlbum(context.Background(), "fooAlbum", "testdata/upload-success")
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		want := mocks.UploadToken + "Id"
		if want != mediaItem.ID {
			t.Errorf("want: %s, got: %s", want, mediaItem.ID)
		}
	})

	t.Run("Should fail with invalid file", func(t *testing.T) {
		client, err := gphotos.NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}

		_, err = client.UploadToAlbum(context.Background(), "fooAlbum", "non-existent")
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})
}
