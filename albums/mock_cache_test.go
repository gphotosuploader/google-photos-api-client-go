package albums

import (
	"context"

	"github.com/duffpl/google-photos-api-client/albums"
)

// Cache mocks a Cache service.
type Cache struct {
	GetAlbumFn      func(ctx context.Context, title string) (albums.Album, error)
	GetAlbumInvoked bool

	PutAlbumFn      func(ctx context.Context, album albums.Album) error
	PutAlbumInvoked bool

	InvalidateAlbumFn      func(ctx context.Context, title string) error
	InvalidateAlbumInvoked bool

	InvalidateAllAlbumsFn      func(ctx context.Context) error
	InvalidateAllAlbumsInvoked bool
}

// GetAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) GetAlbum(ctx context.Context, title string) (albums.Album, error) {
	c.GetAlbumInvoked = true
	return c.GetAlbumFn(ctx, title)
}

// PutAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) PutAlbum(ctx context.Context, album albums.Album) error {
	c.PutAlbumInvoked = true
	return c.PutAlbumFn(ctx, album)
}

// InvalidateAlbum invokes the mock implementation and marks the function as invoked.
func (c *Cache) InvalidateAlbum(ctx context.Context, title string) error {
	c.InvalidateAlbumInvoked = true
	return c.InvalidateAlbumFn(ctx, title)
}

// InvalidateAllAlbums invokes the mock implementation and marks the function as invoked.
func (c *Cache) InvalidateAllAlbums(ctx context.Context) error {
	c.InvalidateAllAlbumsInvoked = true
	return c.InvalidateAllAlbumsFn(ctx)
}
