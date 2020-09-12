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
		t.Fatalf("error was not expected at this point. err: %s", err)
	}

	t.Run("ReturnsErrAlbumNotFoundIfAlbumDoesNotExist", func(t *testing.T) {
		feedAlbumGallery(0)
		_, err := c.FindAlbum(ctx, "nonexistent")
		if !errors.Is(err, gphotos.ErrAlbumNotFound) {
			t.Errorf("error was not expected. want: %s, got: %s", gphotos.ErrAlbumNotFound, err)
		}
	})

	t.Run("ReturnsAlbumIfItIsCached", func(t *testing.T) {
		feedAlbumGallery(0)
		want := "cached-album" // It's already on the cache, see mockedCache.
		got, err := c.FindAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if want != got.Title {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})

	t.Run("ReturnsAlbumEvenItIsNotCached", func(t *testing.T) {
		feedAlbumGallery(5) // Creates five albums: album-1 to album-5.
		want := "album-1"
		got, err := c.FindAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if want != got.Title {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})
}

func TestClient_ListAlbums(t *testing.T) {
	ctx := context.Background()
	c, err := gphotos.NewClient(http.DefaultClient, gphotos.WithPhotoService(mockedService), gphotos.WithCacher(mockedCache))
	if err != nil {
		t.Fatalf("error was not expected at this point. err: %s", err)
	}

	var tests = []struct {
		name string
		want int
	}{
		{name: "Gallery without albums", want: 0},
		{name: "Gallery with on album", want: 1},
		{name: "Gallery with five albums", want: 5},
		{name: "Gallery with fifty albums", want: 50},
		{name: "Gallery with fifty one albums", want: 51},
		{name: "Gallery with five hundred albums", want: 500},
	}

	for _, tc := range tests {
		feedAlbumGallery(tc.want)
		got, err := c.ListAlbums(ctx)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if len(got) != tc.want {
			t.Errorf("case: %s, want: %d, got: %d", tc.name, tc.want, len(got))
		}
	}

}

func TestClient_CreateAlbum(t *testing.T) {
	ctx := context.Background()
	c, err := gphotos.NewClient(http.DefaultClient, gphotos.WithPhotoService(mockedService), gphotos.WithCacher(mockedCache))
	if err != nil {
		t.Fatalf("error was not expected at this point: err=%v", err)
	}

	t.Run("ReturnsAnExistingAlbum", func(t *testing.T) {
		feedAlbumGallery(1) // album-1 has been created on the gallery.
		want := "album-1"
		got, err := c.CreateAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if got.Title != want {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})

	t.Run("ReturnsCreatedAlbum", func(t *testing.T) {
		feedAlbumGallery(0)
		want := "dummy"
		got, err := c.CreateAlbum(ctx, want)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if got.Title != want {
			t.Errorf("want: %s, got: %s", want, got.Title)
		}
	})

	t.Run("ShouldFailDueToAPIError", func(t *testing.T) {
		want := "should-fail"
		got, err := c.CreateAlbum(ctx, want)
		if err == nil {
			t.Fatalf("error was expected at this point. got: %v", got)
		}
	})

}

// albumGallery represents the Album repository to mock Google Photos calls.
var albumGallery []*photoslibrary.Album

// feedAlbumGallery will add the specified number of albums to the Album gallery.
// All the albums follow the template `album-<number>` where `<number>` is an incremental integer.
func feedAlbumGallery(n int) {
	albumGallery = nil
	for i := 1; i <= n; i++ {
		a := photoslibrary.Album{Title: fmt.Sprintf("album-%d", i)}
		albumGallery = append(albumGallery, &a)
	}
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
	GetAlbumFn: func(ctx context.Context, title string) (album photoslibrary.Album, err error) {
		if title == "cached-album" {
			return photoslibrary.Album{Title: "cached-album"}, nil
		}
		return photoslibrary.Album{}, cache.ErrCacheMiss
	},
	PutAlbumFn: func(ctx context.Context, album photoslibrary.Album, ttl time.Duration) error {
		return nil
	},
	InvalidateAlbumFn: func(ctx context.Context, title string) error {
		return nil
	},
}
