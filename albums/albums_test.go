package albums_test

import (
	"context"
	"errors"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should fail without httpClient", func(t *testing.T) {
		cfg := albums.Config{}
		_, err := albums.New(cfg)
		if err == nil {
			t.Errorf("error was expected but not produced")
		}
	})

	t.Run("Should success with an httpClient", func(t *testing.T) {
		cfg := albums.Config{
			Client: http.DefaultClient,
		}
		_, err := albums.New(cfg)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})

	t.Run("Should success with a custom User Agent", func(t *testing.T) {
		cfg := albums.Config{
			Client:    http.DefaultClient,
			UserAgent: "testing-agent",
		}
		_, err := albums.New(cfg)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}

func TestAlbumsService_AddMediaItems(t *testing.T) {
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

	c := albums.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	s, err := albums.New(c)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.AddMediaItems(context.Background(), tc.album, tc.mediaItems)
			assertExpectedError(tc.isErrExpected, err, t)
		})
	}
}

func TestAlbumsService_Create(t *testing.T) {
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

	c := albums.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	s, err := albums.New(c)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Create(context.Background(), tc.input)
			assertExpectedError(tc.isErrExpected, err, t)
			if !tc.isErrExpected && got.ID != tc.want {
				t.Errorf("want: %s, got: %s", tc.want, got.ID)
			}
		})
	}
}

func TestAlbumsService_GetById(t *testing.T) {
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

	c := albums.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	s, err := albums.New(c)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			album, err := s.GetById(context.Background(), tc.input)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("not expected error, want: %v, got: %v", tc.expectedError, err)
			}
			if err == nil && album.ID != tc.input {
				t.Errorf("want: %s, got: %s", tc.input, album.Title)
			}
		})
	}
}

func TestAlbumsService_GetByTitle(t *testing.T) {
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

	c := albums.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	s, err := albums.New(c)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.GetByTitle(context.Background(), tc.input)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("not expected error, want: %v, got: %v", tc.expectedError, err)
			}
			if err == nil && tc.want != got.ID {
				t.Errorf("want: %s, got: %s", tc.want, got.ID)
			}
		})
	}
}

func TestAlbumsService_List(t *testing.T) {
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	c := albums.Config{
		Client:  http.DefaultClient,
		BaseURL: srv.URL(),
	}
	s, err := albums.New(c)
	if err != nil {
		t.Fatalf("error was not expected at this point")
	}

	res, err := s.List(context.Background())
	if err != nil {
		t.Fatal("error was not expected at this point")
	}

	if len(res) != mocks.AvailableAlbums {
		t.Errorf("want: %d, got: %d", mocks.AvailableAlbums, len(res))
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
