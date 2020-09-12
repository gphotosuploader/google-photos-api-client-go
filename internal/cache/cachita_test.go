package cache_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
)

func TestCachitaCache(t *testing.T) {
	c := cache.NewCachitaCache()
	ctx := context.Background()

	// test cache miss
	if _, err := c.GetAlbum(ctx, "nonexistent"); !errors.Is(err, cache.ErrCacheMiss) {
		t.Errorf("want: %v, got: %v", cache.ErrCacheMiss, err)
	}

	// test put/get
	b1 := photoslibrary.Album{Title: "album1"}
	if err := c.PutAlbum(ctx, b1, 60*time.Minute); err != nil {
		t.Fatalf("put: %v", err)
	}
	b2, err := c.GetAlbum(ctx, b1.Title)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !reflect.DeepEqual(b1, b2) {
		t.Errorf("want: %v, got: %v", b1, b2)
	}

	// test delete
	if err := c.InvalidateAlbum(ctx, "dummy"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := c.GetAlbum(ctx, "dummy"); !errors.Is(err, cache.ErrCacheMiss) {
		t.Errorf("want: %v, got: %v", cache.ErrCacheMiss, err)
	}
}
