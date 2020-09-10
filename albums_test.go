package gphotos_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
)

func TestClient_FindAlbum(t *testing.T) {
	ctx := context.Background()
	c, err := gphotos.NewClient(http.DefaultClient, gphotos.WithPhotoService(mockedService), gphotos.WithCacher(mockedCache))
	if err != nil {
		t.Fatalf("error was not expected at this point: err=%v", err)
	}

	t.Run("WithNonExistentAlbum", func(t *testing.T) {
		truncateAlbumGallery()
		_, err := c.FindAlbum(ctx, "nonexistent")
		if !errors.Is(err, gphotos.ErrAlbumNotFound) {
			t.Errorf("error was not expected. want: %v, got: %v", gphotos.ErrAlbumNotFound, err)
		}
	})

	t.Run("WithCachedAlbum", func(t *testing.T) {
		truncateAlbumGallery()
		want := "cached"

		got, err := c.FindAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if want != got.Title {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("WithNonCachedAlbum", func(t *testing.T) {
		initializeAlbumGallery(5)
		want := "album-1"
		got, err := c.FindAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if want != got.Title {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})
}

func TestClient_ListAlbums(t *testing.T) {
	ctx := context.Background()
	c, err := gphotos.NewClient(http.DefaultClient, gphotos.WithPhotoService(mockedService), gphotos.WithCacher(mockedCache))
	if err != nil {
		t.Fatalf("error was not expected at this point: err=%v", err)
	}

	t.Run("WithEmptyAlbumGallery", func(t *testing.T) {
		truncateAlbumGallery()
		got, err := c.ListAlbums(ctx)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if len(got) > 0 {
			t.Errorf("no albums should be listed. got: %d", len(got))
		}
	})

	t.Run("WithSmallAlbumGallery", func(t *testing.T) {
		want := 5
		initializeAlbumGallery(want)
		got, err := c.ListAlbums(ctx)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if len(got) != want {
			t.Errorf("want: %d, got: %d", want, len(got))
		}
	})

	t.Run("WithLargeAlbumGallery", func(t *testing.T) {
		want := 500
		initializeAlbumGallery(want)
		got, err := c.ListAlbums(ctx)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if len(got) != want {
			t.Errorf("want: %d, got: %d", want, len(got))
		}
	})

}

func TestClient_CreateAlbum(t *testing.T) {
	ctx := context.Background()
	c, err := gphotos.NewClient(http.DefaultClient, gphotos.WithPhotoService(mockedService), gphotos.WithCacher(mockedCache))
	if err != nil {
		t.Fatalf("error was not expected at this point: err=%v", err)
	}

	t.Run("ReturnsExistingAlbum", func(t *testing.T) {
		initializeAlbumGallery(1)
		want := "album-1"
		got, err := c.CreateAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if got.Title != want {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})

	t.Run("ReturnsCreatedAlbum", func(t *testing.T) {
		truncateAlbumGallery()
		want := "dummy"
		got, err := c.CreateAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %v", err)
		}
		if got.Title != want {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})

	t.Run("ShouldFailDueToAPIError", func(t *testing.T) {
		truncateAlbumGallery()
		want := "should-fail"
		got, err := c.CreateAlbum(ctx, want)
		if err == nil {
			t.Fatalf("error was expected at this point. got: %v", got)
		}
	})

}

// albumGallery represents the Album repository to mock Google Photos calls.
var albumGallery []*photoslibrary.Album

// initializeAlbumGallery will add the specified number of albums to the Album gallery.
// All the albums follow the template `album-<number>` where `<number>` is an incremental integer.
func initializeAlbumGallery(n int) {
	truncateAlbumGallery()
	for i := 1; i <= n; i++ {
		a := photoslibrary.Album{Title: fmt.Sprintf("album-%d", i)}
		albumGallery = append(albumGallery, &a)
	}
}

// truncateAlbumGallery will empty the Album gallery.
func truncateAlbumGallery() {
	albumGallery = nil
}

var mockedService = &mock.PhotoService{
	ListAlbumsFn: func(ctx context.Context, pageSize int64, pageToken string) (response *photoslibrary.ListAlbumsResponse, err error) {
		if pageToken == "give-me-more" {
			// second page of albums.
			return &photoslibrary.ListAlbumsResponse{
				Albums: albumGallery[pageSize:],
			}, nil
		}

		if pageSize < int64(len(albumGallery)) {
			// first page of albums.
			return &photoslibrary.ListAlbumsResponse{
				Albums:        albumGallery[:pageSize],
				NextPageToken: "give-me-more",
			}, nil
		}

		// there is only one page of albums.
		return &photoslibrary.ListAlbumsResponse{
			Albums: albumGallery,
		}, nil
	},
	CreateAlbumFn: func(ctx context.Context, request *photoslibrary.CreateAlbumRequest) (album *photoslibrary.Album, err error) {
		if request.Album.Title == "should-fail" {
			return nil, errors.New("album creation failure")
		}
		return request.Album, nil
	},
	CreateMediaItemsFn: func(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (response *photoslibrary.BatchCreateMediaItemsResponse, err error) {
		return &photoslibrary.BatchCreateMediaItemsResponse{}, nil
	},
}

var mockedCache = &mock.Cache{
	GetAlbumFn: func(ctx context.Context, title string) (album *photoslibrary.Album, err error) {
		if title == "cached" {
			return &photoslibrary.Album{Title: "cached"}, nil
		}
		return nil, cache.ErrCacheMiss
	},
	PutAlbumFn: func(ctx context.Context, album *photoslibrary.Album, ttl time.Duration) error {
		return nil
	},
	InvalidateAlbumFn: func(ctx context.Context, title string) error {
		return nil
	},
}
