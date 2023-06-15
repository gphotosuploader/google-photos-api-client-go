package gphotos_test

import (
	"context"
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"net/http"
	"testing"
)

func TestClient_UploadFileToLibrary(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if upload fails", "upload-should-fail", true},
		{"Should return error if mediaItemManager fails", "media-item-manager-should-fail", true},
		{"Should return success on success", "foo", false},
	}

	mockedUploader := &mocks.MockedUploader{
		UploadFileFn: func(ctx context.Context, filePath string) (uploadToken string, err error) {
			if "upload-should-fail" == filePath {
				return "", fmt.Errorf("uploader-error")
			}
			return "token", nil
		},
	}

	mockedMediaItemManager := &mocks.MockedMediaItemsService{
		CreateFn: func(ctx context.Context, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error) {
			if "media-item-manager-should-fail" == mediaItem.FileName {
				return media_items.MediaItem{}, fmt.Errorf("media-item-manager-error")
			}
			return media_items.MediaItem{Filename: mediaItem.FileName}, nil
		}}

	config := gphotos.Config{
		Client:           http.DefaultClient,
		Uploader:         mockedUploader,
		MediaItemManager: mockedMediaItemManager,
	}
	client, err := gphotos.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := client.UploadFileToLibrary(context.Background(), tc.input)
			if tc.isErrExpected && err == nil {
				t.Fatalf("error was expected, but not produced")
			}
			if !tc.isErrExpected && got.Filename != tc.input {
				t.Errorf("want: %s, got: %s", tc.input, got.Filename)
			}
		})
	}
}

func TestClient_UploadFileToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if upload fails", "upload-should-fail", true},
		{"Should return error if mediaItemManager fails", "media-item-manager-should-fail", true},
		{"Should return success on success", "foo", false},
	}

	mockedUploader := &mocks.MockedUploader{
		UploadFileFn: func(ctx context.Context, filePath string) (uploadToken string, err error) {
			if "upload-should-fail" == filePath {
				return "", fmt.Errorf("uploader-error")
			}
			return "token", nil
		},
	}

	mockedMediaItemManager := &mocks.MockedMediaItemsService{
		CreateToAlbumFn: func(ctx context.Context, albumId string, mediaItem media_items.SimpleMediaItem) (media_items.MediaItem, error) {
			if "media-item-manager-should-fail" == mediaItem.FileName {
				return media_items.MediaItem{}, fmt.Errorf("media-item-manager-error")
			}
			return media_items.MediaItem{Filename: mediaItem.FileName}, nil
		},
	}

	config := gphotos.Config{
		Client:           http.DefaultClient,
		Uploader:         mockedUploader,
		MediaItemManager: mockedMediaItemManager,
	}
	client, err := gphotos.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := client.UploadFileToAlbum(context.Background(), "album", tc.input)
			if tc.isErrExpected && err == nil {
				t.Fatalf("error was expected, but not produced")
			}
			if !tc.isErrExpected && got.Filename != tc.input {
				t.Errorf("want: %s, got: %s", tc.input, got.Filename)
			}
		})
	}
}
