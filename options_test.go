package gphotos_test

import (
	"context"
	"testing"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/mock"
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
	want := &mockedCache{}

	got := gphotos.WithCacher(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

type mockedCache struct{}

func (mc *mockedCache) GetAlbum(ctx context.Context, key string) (*photoslibrary.Album, error) {
	if key == "cached" {
		return &photoslibrary.Album{Title: "cached"}, nil
	}
	return nil, cache.ErrCacheMiss
}

func (mc *mockedCache) PutAlbum(ctx context.Context, key string, album *photoslibrary.Album, ttl time.Duration) error {
	return nil
}
func (mc *mockedCache) InvalidateAlbum(ctx context.Context, key string) error {
	return nil
}
