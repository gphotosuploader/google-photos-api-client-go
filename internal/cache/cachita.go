package cache

import (
	"context"
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// CachitaCache implements Cache with `gadelkareem/cachita` package.
type CachitaCache struct {
	store cachita.Cache
}

// NewCachitaCache returns a Cache service implemented using `gadelkareem/cachita`.
func NewCachitaCache() *CachitaCache {
	return &CachitaCache{store: cachita.Memory()}
}

// Get reads an object data from the cache.
func (c *CachitaCache) GetAlbum(ctx context.Context, title string) (photoslibrary.Album, error) {
	item := photoslibrary.Album{}
	err := c.store.Get(c.albumKey(title), &item)
	if err == cachita.ErrNotFound {
		return photoslibrary.Album{}, ErrCacheMiss
	}

	return item, err
}

// Put store an object data to the cache.
func (c *CachitaCache) PutAlbum(ctx context.Context, album photoslibrary.Album, ttl time.Duration) error {
	return c.store.Put(c.albumKey(album.Title), album, ttl)
}

// InvalidateAlbum removes the specified Album from the cache.
func (c *CachitaCache) InvalidateAlbum(ctx context.Context, title string) error {
	return c.store.Invalidate(c.albumKey(title))
}

// albumKey returns the cache key for an Album title.
func (c *CachitaCache) albumKey(title string) string {
	return "album:" + title
}
