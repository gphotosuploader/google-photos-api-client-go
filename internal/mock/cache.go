package mock

import (
	"context"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

type Cache struct {
	GetAlbumFn      func(ctx context.Context, title string) (photoslibrary.Album, error)
	GetAlbumInvoked bool

	PutAlbumFn      func(ctx context.Context, album photoslibrary.Album, ttl time.Duration) error
	PutAlbumInvoked bool

	InvalidateAlbumFn      func(ctx context.Context, title string) error
	InvalidateAlbumInvoked bool
}

// GetAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) GetAlbum(ctx context.Context, title string) (photoslibrary.Album, error) {
	c.GetAlbumInvoked = true
	return c.GetAlbumFn(ctx, title)
}

// PutAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) PutAlbum(ctx context.Context, album photoslibrary.Album, ttl time.Duration) error {
	c.PutAlbumInvoked = true
	return c.PutAlbumFn(ctx, album, ttl)
}

// InvalidateAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) InvalidateAlbum(ctx context.Context, title string) error {
	c.InvalidateAlbumInvoked = true
	return c.InvalidateAlbumFn(ctx, title)
}
