package cache_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/duffpl/google-photos-api-client/albums"

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
	b1 := albums.Album{Title: "album1"}
	if err := c.PutAlbum(ctx, b1); err != nil {
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

func TestCachitaCache_PutAlbum(t *testing.T) {
	testCases := []struct {
		name        string
		input       albums.Album
		errExpected bool
	}{
		{"empty album", albums.Album{}, false},
		{"album with title", albums.Album{Title: "foo"}, false},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.NewCachitaCache()
			err := c.PutAlbum(ctx, tc.input)
			if tc.errExpected && err == nil {
				t.Errorf("error was expected, but not produced")
			}
			if !tc.errExpected && err != nil {
				t.Errorf("error was not expected. err: %s", err)
			}
		})
	}
}

func TestCachitaCache_GetAlbum(t *testing.T) {
	testCases := []struct {
		name           string
		populatedCache []string
		input          string
		errExpected    error
	}{
		{"empty cache", []string{}, "foo", cache.ErrCacheMiss},
		{"existing key", []string{"foo", "bar"}, "foo", nil},
		{"non-existent key", []string{"foo", "bar"}, "baz", cache.ErrCacheMiss},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.NewCachitaCache()
			for _, title := range tc.populatedCache {
				if err := c.PutAlbum(ctx, albums.Album{Title: title}); err != nil {
					t.Fatalf("error was not expected at this point. err: %s", err)
				}
			}
			_, err := c.GetAlbum(ctx, tc.input)
			if tc.errExpected != err {
				t.Errorf("not expected error, want: %v, got: %v", tc.errExpected, err)
			}
		})
	}
}

func TestCachitaCache_InvalidateAlbum(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		errExpected bool
	}{
		{"existing key", "foo", false},
		{"non-existent key", "dummy", false},
	}

	ctx := context.Background()
	populatedCache := []string{"foo", "bar", "baz"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := cache.NewCachitaCache()
			for _, title := range populatedCache {
				if err := c.PutAlbum(ctx, albums.Album{Title: title}); err != nil {
					t.Fatalf("error was not expected at this point. err: %s", err)
				}
			}
			err := c.InvalidateAlbum(ctx, tc.input)
			if tc.errExpected && err == nil {
				t.Errorf("error was expected, but not produced")
			}
			if !tc.errExpected && err != nil {
				t.Errorf("error was not expected. err: %s", err)
			}
		})
	}
}
