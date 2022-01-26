package albums_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
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
			_, err := albums.NewPhotosLibraryClient(tc.input)
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
			ar, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && ar.URL() != tc.want {
				t.Errorf("want: %s, got: %s", tc.want, ar.URL())
			}
		})
	}
}

func TestPhotosLibraryAlbumsRepository_AddManyItems(t *testing.T) {
	testCases := []struct {
		name          string
		album         string
		mediaItems    []string
		isErrExpected bool
	}{
		{"Should add media items to album", "foo", []string{"mediaItem1", "mediaItem2"}, true},
		{"Should return error if album does not exist", "non-existent", []string{"mediaItem1", "mediaItem2"}, true},
		{"Should return error if media item is invalid", "foo", []string{mocks.ShouldMakeAPIFailMediaItem, "mediaItem2"}, true},
		{"Should return error if API fails", mocks.ShouldFailAlbum.Id, []string{"mediaItem1", "mediaItem2"}, true},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	r, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := r.AddManyItems(context.Background(), tc.album, tc.mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestPhotosLibraryAlbumsRepository_Create(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should return error if API fails", mocks.ShouldFailAlbum.Title, "", true},
		{"Should return the album on success", "foo", "fooId", false},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	ar, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			album, err := ar.Create(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && album.ID != tc.want {
				t.Errorf("want: %s, got: %s", tc.want, album.ID)
			}
		})
	}
}

func TestPhotosLibraryAlbumsRepository_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{"Should return the album on success", "fooId-0", nil},
		{"Should return ErrAlbumNotFound if API fails", mocks.ShouldFailAlbum.Id, albums.ErrAlbumNotFound},
		{"Should return ErrAlbumNotFound if albums does not exist", "non-existent", albums.ErrAlbumNotFound},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	ar, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			album, err := ar.Get(context.Background(), tc.input)
			if tc.expectedError != err {
				t.Fatalf("not expected error, want: %v, got: %v", tc.expectedError, err)
			}
			if err == nil && album.ID != tc.input {
				t.Errorf("want: %s, got: %s", tc.input, album.Title)
			}
		})
	}
}

func TestPhotosLibraryAlbumsRepository_GetByTitle(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		expectedError error
	}{
		{"Should return the album on success", "fooTitle-0", "fooId-0", nil},
		{"Should return ErrAlbumNotFound if API fails", mocks.ShouldFailAlbum.Id, "", albums.ErrAlbumNotFound},
		{"Should return ErrAlbumNotFound if the album does not exist", "non-existent", "", albums.ErrAlbumNotFound},
	}

	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	ar, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ar.GetByTitle(context.Background(), tc.input)
			if tc.expectedError != err {
				t.Fatalf("not expected error, want: %v, got: %v", tc.expectedError, err)
			}
			if err == nil && tc.want != got.ID {
				t.Errorf("want: %s, got: %s", tc.want, got.ID)
			}
		})
	}
}

func TestPhotosLibraryAlbumsRepository_ListAll(t *testing.T) {
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	albumsService, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	res, err := albumsService.ListAll(context.Background())
	if err != nil {
		t.Fatal("error was not expected at this point")
	}

	if len(res) != mocks.AvailableAlbums {
		t.Errorf("want: %d, got: %d", mocks.AvailableAlbums, len(res))
	}
}

func TestPhotosLibraryAlbumsRepository_ListWithOptions(t *testing.T) {
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	albumsService, err := albums.NewPhotosLibraryClientWithURL(http.DefaultClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	res, err := albumsService.ListWithOptions(context.Background(), albums.Options{ExcludeNonAppCreatedData: true})
	if err != nil {
		t.Fatal("error was not expected at this point")
	}

	if len(res) != mocks.AvailableAlbums {
		t.Errorf("want: %d, got: %d", mocks.AvailableAlbums, len(res))
	}
}
