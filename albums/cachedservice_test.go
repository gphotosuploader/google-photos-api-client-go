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

func TestCachedAlbumsService_Create(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))

	t.Run("Should return error if API fails", func(t *testing.T) {
		_, err := s.Create("api-should-fail", context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return error if cache fails", func(t *testing.T) {
		_, err := s.Create("cache-should-fail", context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return the created album on success", func(t *testing.T) {
		_, err := s.Create("foo", context.Background())
		if err != nil {
			t.Errorf("error was not expected, err: %s", err)
		}
	})
}

func TestCachedAlbumsService_Get(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))

	t.Run("Should return error if API fails", func(t *testing.T) {
		_, err := s.Get("api-should-fail", context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return error if cache fails", func(t *testing.T) {
		_, err := s.Get("cache-should-fail", context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return the album on success", func(t *testing.T) {
		_, err := s.Get("foo", context.Background())
		if err != nil {
			t.Errorf("error was not expected, err: %s", err)
		}
	})
}

func TestCachedAlbumsService_GetByTitle(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))

	t.Run("Should return error if cache fails", func(t *testing.T) {
		_, err := s.GetByTitle("cache-should-fail", context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return the album on success", func(t *testing.T) {
		album, err := s.GetByTitle("foo", context.Background())
		if err != nil {
			t.Fatalf("error not expected at this point, err: %s", err)
		}
		if "foo" != album.Title {
			t.Errorf("wamt: %s, got: %s", "foo", album.Title)
		}
	})

	t.Run("Should return ErrAlbumNotFound if the album does not exist", func(t *testing.T) {
		_, err := s.GetByTitle("non-existent", context.Background())
		if err != ErrAlbumNotFound {
			t.Errorf("error was not expected, err: %s", err)
		}
	})
}

func TestCachedAlbumsService_ListAll(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))

	t.Run("Should return the existing albums", func(t *testing.T) {
		res, err := s.ListAll(&albums.AlbumsListOptions{}, context.Background())
		if err != nil {
			t.Fatalf("error was expected at this point")
		}
		if 3 != len(res) {
			t.Errorf("#albums, want: %d, got: %d", 3, len(res))
		}
	})
}

func TestCachedAlbumsService_Patch(t *testing.T) {
	s := NewCachedAlbumsService(http.DefaultClient, WithAlbumsAPIClient(mockedAlbumsAPIClient), WithCacher(mockedCache))

	t.Run("Should return error if API fails", func(t *testing.T) {
		album := albums.Album{Title:"api-should-fail"}
		updateMask := []albums.Field{albums.AlbumFieldTitle}
		_, err := s.Patch(album, updateMask, context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should  if API fails", func(t *testing.T) {
		album := albums.Album{Title:"api-should-fail"}
		updateMask := []albums.Field{albums.AlbumFieldTitle}
		_, err := s.Patch(album, updateMask, context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return error if cache fails", func(t *testing.T) {
		album := albums.Album{Title:"cache-should-fail"}
		updateMask := []albums.Field{albums.AlbumFieldTitle}
		_, err := s.Patch(album, updateMask, context.Background())
		if err == nil {
			t.Errorf("error was expected, but not produced")
		}
	})

	t.Run("Should return an album on success", func(t *testing.T) {
		album := albums.Album{Title:"foo"}
		updateMask := []albums.Field{albums.AlbumFieldTitle}
		got, err := s.Patch(album, updateMask, context.Background())
		if err != nil {
			t.Fatalf("error was expected at this point")
		}
		if album.Title != got.Title {
			t.Errorf("want: %s, got: %s", album.Title, got.Title)
		}
	})
}

var mockedAlbumsAPIClient = mock.AlbumService{
	CreateFn: func(title string, ctx context.Context) (*albums.Album, error) {
		if title == "api-should-fail" {
			return &albums.Album{}, errors.New("error")
		}
		return &albums.Album{Title: title}, nil
	},
	GetFn: func(id string, ctx context.Context) (*albums.Album, error) {
		if id == "api-should-fail" {
			return &albums.Album{}, errors.New("error")
		}
		return &albums.Album{ID: id}, nil
	},
	ListAllAsyncFn: func(options *albums.AlbumsListOptions, ctx context.Context) (<-chan albums.Album, <-chan error) {
		items := []string{"foo", "bar", "baz"}

		albumsC := make(chan albums.Album, len(items))
		errorsC := make(chan error)
		go func() {
			defer close(albumsC)
			for _, item := range items {
				albumsC <- albums.Album{Title: item}
			}
		}()
		return albumsC, errorsC
	},
	PatchFn: func(album albums.Album, updateMask []albums.Field, ctx context.Context) (*albums.Album, error) {
		if album.Title == "api-should-fail" {
			return &albums.Album{}, errors.New("error")
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

var mockedCache = &mock.Cache{
	GetAlbumFn: func(ctx context.Context, title string) (album albums.Album, err error) {
		if title == "cached-album" {
			return albums.Album{Title: "cached-album"}, nil
		}
		return albums.Album{}, cache.ErrCacheMiss
	},
	PutAlbumFn: func(ctx context.Context, album albums.Album) error {
		if album.Title == "cache-should-fail" || album.ID == "cache-should-fail" {
			return errors.New("error")
		}
		return nil
	},
	InvalidateAlbumFn: func(ctx context.Context, title string) error {
		return nil
	},
	InvalidateAllAlbumsFn: func(ctx context.Context) error {
		return nil
	},
}
