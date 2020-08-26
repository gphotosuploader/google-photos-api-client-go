package cache_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	cache2 "github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
)

func TestCachitaCache(t *testing.T) {
	cache := cache2.NewCachitaCache()
	ctx := context.Background()

	// test cache miss
	if _, err := cache.GetAlbum(ctx, "nonexistent"); err != cache2.ErrCacheMiss {
		t.Errorf("want: %v, got: %v", cache2.ErrCacheMiss, err)
	}

	// test put/get
	b1 := &photoslibrary.Album{Title: "album1"}
	if err := cache.PutAlbum(ctx, "dummy", b1, 60 * time.Minute); err != nil {
		t.Fatalf("put: %v", err)
	}
	b2, err := cache.GetAlbum(ctx, "dummy")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !reflect.DeepEqual(b1, b2) {
		t.Errorf("want: %v, got: %v", b1, b2)
	}

	// test delete
	if err := cache.InvalidateAlbum(ctx, "dummy"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := cache.GetAlbum(ctx, "dummy"); err != cache2.ErrCacheMiss {
		t.Errorf("want: %v, got: %v", cache2.ErrCacheMiss, err)
	}
}
