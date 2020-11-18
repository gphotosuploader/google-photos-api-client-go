package albums

import (
	"context"
)

// MockedCache mocks a Cache service.
type MockedCache struct {
	GetAlbumFn            func(ctx context.Context, title string) (Album, error)
	PutAlbumFn            func(ctx context.Context, album Album) error
	InvalidateAlbumFn     func(ctx context.Context, title string) error
	InvalidateAllAlbumsFn func(ctx context.Context) error
}

// GetAlbum invokes the mock implementation and marks the function as invoked.
func (c MockedCache) GetAlbum(ctx context.Context, title string) (Album, error) {
	return c.GetAlbumFn(ctx, title)
}

// PutAlbum invokes the mock implementation and marks the function as invoked.
func (c MockedCache) PutAlbum(ctx context.Context, album Album) error {
	return c.PutAlbumFn(ctx, album)
}

// InvalidateAlbum invokes the mock implementation and marks the function as invoked.
func (c MockedCache) InvalidateAlbum(ctx context.Context, title string) error {
	return c.InvalidateAlbumFn(ctx, title)
}

// InvalidateAllAlbums invokes the mock implementation and marks the function as invoked.
func (c MockedCache) InvalidateAllAlbums(ctx context.Context) error {
	return c.InvalidateAllAlbumsFn(ctx)
}
