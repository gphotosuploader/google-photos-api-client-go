package media_items_test

import (
	"context"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"net/http"
	"testing"
)

func TestHttpMediaItemsService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		uploadToken   string
		isErrExpected bool
	}{
		{"Should return error if API fails", mocks.ShouldMakeAPIFailMediaItem, true},
		{"Should return success on success", "foo", false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := m.Create(context.Background(), media_items.SimpleMediaItem{UploadToken: tc.uploadToken})
			assertExpectedError(tc.isErrExpected, err, t)
			want := tc.uploadToken + "Id"
			if err == nil && want != got.ID {
				t.Errorf("want: %s, got: %s", want, got.ID)
			}
		})
	}
}

func TestHttpMediaItemsService_CreateMany(t *testing.T) {
	testCases := []struct {
		name          string
		uploadTokens  []string
		want          int
		isErrExpected bool
	}{
		{"Should return error if API fails", []string{mocks.ShouldMakeAPIFailMediaItem, "dummy"}, 0, true},
		{"Should return success on success", []string{"foo", "bar", "baz"}, 3, false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := make([]media_items.SimpleMediaItem, len(tc.uploadTokens))
			for i, token := range tc.uploadTokens {
				mediaItems[i].UploadToken = token
			}
			got, err := m.CreateMany(context.Background(), mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != len(got) {
				t.Errorf("want: %d, got: %d", tc.want, len(got))
			}
		})
	}
}

func TestHttpMediaItemsService_CreateToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		albumId       string
		uploadToken   string
		isErrExpected bool
	}{
		{"Should return error if API fails", "albumId", mocks.ShouldMakeAPIFailMediaItem, true},
		{"Should return success on success", "foo", "bar", false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := m.CreateToAlbum(context.Background(), tc.albumId, media_items.SimpleMediaItem{UploadToken: tc.uploadToken})
			assertExpectedError(tc.isErrExpected, err, t)
			want := tc.uploadToken + "Id"
			if err == nil && want != got.ID {
				t.Errorf("want: %s, got: %s", want, got.ID)
			}
		})
	}
}

func TestHttpMediaItemsService_CreateManyToAlbum(t *testing.T) {
	testCases := []struct {
		name          string
		albumId       string
		uploadTokens  []string
		want          int
		isErrExpected bool
	}{
		{"Should return error if API fails", "foo", []string{mocks.ShouldMakeAPIFailMediaItem, "dummy-2"}, 0, true},
		{"Should return success on success", "foo", []string{"bar", "baz"}, 2, false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mediaItems := make([]media_items.SimpleMediaItem, len(tc.uploadTokens))
			for i, token := range tc.uploadTokens {
				mediaItems[i].UploadToken = token
			}
			got, err := m.CreateManyToAlbum(context.Background(), tc.albumId, mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != len(got) {
				t.Errorf("want: %d, got: %d", tc.want, len(got))
			}
		})
	}
}

func TestHttpMediaItemsService_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should return error if API fails", mocks.ShouldMakeAPIFailMediaItem, "", true},
		{"Should return error if media item does not exist", "non-existent", "", true},
		{"Should return success on success", "fooId-0", "fooId-0", false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := m.Get(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.want != got.ID {
				t.Errorf("want: %s, got: %s", tc.want, got.ID)
			}
		})
	}
}

func TestHttpMediaItemsService_ListByAlbum(t *testing.T) {
	testCases := []struct {
		name  string
		input string

		isErrExpected bool
	}{
		{"Should return error if API fails", mocks.ShouldFailAlbum.Id, true},
		{"Should return success on success", "fooId-0", false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	config := media_items.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	m, err := media_items.New(config)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := m.ListByAlbum(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && mocks.AvailableMediaItems != len(got) {
				t.Errorf("want: %d, got: %d", mocks.AvailableMediaItems, len(got))
			}
		})
	}
}

func assertExpectedError(isErrExpected bool, err error, t *testing.T) {
	if isErrExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !isErrExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
