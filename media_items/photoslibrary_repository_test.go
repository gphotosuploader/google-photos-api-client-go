package media_items_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mocks"
)

const defaultBasePath = "https://photoslibrary.googleapis.com/"

func TestNewPhotosLibraryClient(t *testing.T) {
	testCases := []struct {
		name          string
		input         *http.Client
		isErrExpected bool
	}{
		{"No HTTP client", nil, true},
		{"Default HTTP client", http.DefaultClient, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := media_items.NewPhotosLibraryClient(tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestNewPhotosLibraryClientWithURL(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"New client with defaults", "", defaultBasePath, false},
		{"New client with custom URL", "https://mydomain.com/path/", "https://mydomain.com/path/", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ar, err := media_items.NewPhotosLibraryClientWithURL(http.DefaultClient, tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && ar.URL() != tc.want {
				t.Errorf("want: %s, got: %s", tc.want, ar.URL())
			}
		})
	}
}

func TestPhotosLibraryMediaItemsRepository_CreateMany(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		isErrExpected bool
	}{
		{"Should return error if API fails", []string{mocks.ShouldMakeAPIFailMediaItem}, true},
		{"Should return a media items on success", []string{"foo"}, false},
		{"Should return multiple media items on success", []string{"foo", "bar"}, false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	r, err := media_items.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := createMediaItems(tc.input)
			result, err := r.CreateMany(context.Background(), mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && len(tc.input) != len(result) {
				t.Errorf("want: %d, got: %d", len(result), len(tc.input))
			}
		})
	}
}

func TestPhotosLibraryMediaItemsRepository_CreateManyToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		album         string
		input         []string
		isErrExpected bool
	}{
		{"Should return error if API fails", "albumId", []string{mocks.ShouldMakeAPIFailMediaItem}, true},
		{"Should return a media items on success", "albumId", []string{"foo"}, false},
		{"Should return multiple media items on success", "albumId", []string{"foo", "bar"}, false},
		{"Should return when one media item fails (issue #54)", "", []string{"foo", mocks.ShouldReturnEmptyMediaItem, "bar"}, false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	r, err := media_items.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := createMediaItems(tc.input)
			result, err := r.CreateManyToAlbum(context.Background(), "albumId", mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && len(tc.input) != len(result) {
				t.Errorf("want: %d, got: %d", len(result), len(tc.input))
			}
		})
	}
}

func TestPhotosLibraryMediaItemsRepository_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		expectedError error
	}{
		{"Should return the media item on success", "fooId-0", "fooFilename-0", nil},
		{"Should return ErrServerFailed if API fails", mocks.ShouldMakeAPIFailMediaItem, "", media_items.ErrServerFailed},
		{"Should return ErrNotFound if media item does not exist", "non-existent", "", media_items.ErrNotFound},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	r, err := media_items.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItem, err := r.Get(context.Background(), tc.input)
			if tc.expectedError != err {
				t.Fatalf("not expected error, want: %v, got: %v", tc.expectedError, err)
			}
			if err == nil && mediaItem.Filename != tc.want {
				t.Errorf("want: %s, got: %s", tc.want, mediaItem.Filename)
			}
		})
	}
}

func TestPhotosLibraryMediaItemsRepository_ListByAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return all media items in album", "foo", false},
		{"Should return all media items with pagination", mocks.ShouldReturnPaginatedAlbum.Id, false},
		{"Should return error if API fails", mocks.ShouldFailAlbum.Id, true},
	}
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	r, err := media_items.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems, err := r.ListByAlbum(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && mocks.AvailableMediaItems != len(mediaItems) {
				t.Errorf("want: %d, got: %d", mocks.AvailableMediaItems, len(mediaItems))
			}
		})
	}
}

func createMediaItems(input []string) []media_items.SimpleMediaItem {
	mediaItems := make([]media_items.SimpleMediaItem, len(input))
	for i, mi := range input {
		mediaItems[i] = media_items.SimpleMediaItem{
			UploadToken: mi,
		}
	}
	return mediaItems
}

func assertExpectedError(isErrExpected bool, err error, t *testing.T) {
	if isErrExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !isErrExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
