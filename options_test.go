package gphotos_test

import (
	"context"
	"testing"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
)

func TestWithPhotoService(t *testing.T) {
	want := &mock.PhotoService{}

	got := gphotos.WithPhotoService(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithUploader(t *testing.T) {
	want := &mock.Uploader{}

	got := gphotos.WithUploader(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithLogger(t *testing.T) {
	want := &log.DiscardLogger{}

	got := gphotos.WithLogger(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithSessionStorer(t *testing.T) {
	want := &mock.SessionStorer{}

	got := gphotos.WithSessionStorer(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithCacher(t *testing.T) {
	want := &mock.Cache{
		GetAlbumFn: func(ctx context.Context, title string) (album photoslibrary.Album, err error) {
			if title == "cached" {
				return photoslibrary.Album{Title: "cached"}, nil
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

	got := gphotos.WithCacher(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
