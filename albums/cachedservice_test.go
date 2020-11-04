package albums

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/duffpl/google-photos-api-client/albums"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
)

var (
	mockedAlbumsRepository = make([]albums.Album, 0)

	mockedAlbumsAPIClient = mock.AlbumService{
		CreateFn: func(title string, ctx context.Context) (*albums.Album, error) {
			if title == "create-api-should-fail" {
				return &NullAlbum, errors.New("error")
			}
			return &albums.Album{Title: title}, nil
		},
		GetFn: func(id string, ctx context.Context) (*albums.Album, error) {
			if id == "get-api-should-fail" {
				return &NullAlbum, errors.New("error")
			}
			return &albums.Album{ID: id}, nil
		},
		ListAllAsyncFn: func(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error) {
			albumsC := make(chan albums.Album, len(mockedAlbumsRepository))
			errorsC := make(chan error)
			go func() {
				defer close(albumsC)
				for _, item := range mockedAlbumsRepository {
					albumsC <- item
				}
			}()
			return albumsC, errorsC
		},
		PatchFn: func(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error) {
			if album.Title == "patch-api-should-fail" {
				return &NullAlbum, errors.New("error")
			}
			original := albums.Album{
				ID:                    "originalAlbumId",
				Title:                 "originalTitle",
				CoverPhotoMediaItemID: "originalCoverPhotoMediaItemId",
			}
			for _, field := range updateMask {
				switch field {
				case albums.AlbumFieldTitle:
					original.Title = album.Title
				case albums.AlbumFieldCoverPhotoMediaItemId:
					original.CoverPhotoMediaItemID = album.CoverPhotoMediaItemID
				}
			}
			return &original, nil
		},
	}

	mockedCache = &mock.Cache{
		GetAlbumFn: func(ctx context.Context, title string) (album albums.Album, err error) {
			if title == "get-cache-should-fail" {
				return NullAlbum, errors.New("error")
			}
			if title == "cached-album" {
				return albums.Album{Title: "cached-album"}, nil
			}
			return NullAlbum, cache.ErrCacheMiss
		},
		PutAlbumFn: func(ctx context.Context, album albums.Album) error {
			if album.Title == "put-cache-should-fail" || album.ID == "put-cache-should-fail" {
				return errors.New("error")
			}
			return nil
		},
		InvalidateAlbumFn: func(ctx context.Context, title string) error {
			if title == "invalidate-cache-should-fail" {
				return errors.New("error")
			}
			return nil
		},
		InvalidateAllAlbumsFn: func(ctx context.Context) error {
			return nil
		},
	}
)

func TestCachedAlbumsService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "create-api-should-fail", true},
		{"Should return error if cache fails", "put-cache-should-fail", true},
		{"Should return the created album on success", "foo", false},
	}
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Create(tc.input, context.Background())
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestCachedAlbumsService_Get(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "get-api-should-fail", true},
		{"Should return error if cache fails", "put-cache-should-fail", true},
		{"Should return the created album on success", "foo", false},
	}
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Get(tc.input, context.Background())
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.ID {
				t.Errorf("want: %s, got: %s", tc.input, got.ID)
			}
		})
	}
}

func TestCachedAlbumsService_GetByTitle(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
		errExpected   error
	}{
		{"Should return error if cache fails (get)", "get-cache-should-fail", true, nil},
		{"Should return the cached album on success", "cached-album", false, nil},
		{"Should return the album on success", "bar", false, nil},
		{"Should return ErrAlbumNotFound if the album does not exist", "non-existent", true, ErrAlbumNotFound},
	}
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))
	fillMockedRepository([]string{"foo", "bar", "baz"})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.GetByTitle(tc.input, context.Background())
			assertExpectedError(tc.isErrExpected, err, t)
			if tc.errExpected != nil && tc.errExpected != err {
				t.Errorf("err want: %s, err got: %s", tc.errExpected, err)
			}
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func TestCachedAlbumsService_ListAll(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))
	fillMockedRepository([]string{"foo", "bar", "baz"})

	t.Run("Should return the existing albums", func(t *testing.T) {
		res, err := s.ListAll(&albums.AlbumsListOptions{}, context.Background())
		if err != nil {
			t.Fatalf("error was expected at this point")
		}
		if len(mockedAlbumsRepository) != len(res) {
			t.Errorf("#albums, want: %d, got: %d", 3, len(res))
		}
	})
}

func TestCachedAlbumsService_Patch(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should return error if API fails", "patch-api-should-fail", true},
		{"Should return error if cache fails (invalidate)", "invalidate-cache-should-fail", true},
		{"Should return error if cache fails (put)", "put-cache-should-fail", true},
		{"Should return the modified album on success", "foo", false},
	}
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			album := albums.Album{Title: tc.input}
			updateMask := []albums.Field{albums.AlbumFieldTitle}
			got, err := s.Patch(album, updateMask, context.Background())
			assertExpectedError(tc.isErrExpected, err, t)
			if err == nil && tc.input != got.Title {
				t.Errorf("want: %s, got: %s", tc.input, got.Title)
			}
		})
	}
}

func fillMockedRepository(items []string) {
	mockedAlbumsRepository = make([]albums.Album, 0)
	for _, item := range items {
		mockedAlbumsRepository = append(mockedAlbumsRepository, albums.Album{
			ID:    item + "Id",
			Title: item,
		})
	}
}

func assertExpectedError(errExpected bool, err error, t *testing.T) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
