package cache

import (
	"context"
	"time"

	"github.com/duffpl/google-photos-api-client/albums"
	"github.com/gadelkareem/cachita"
)

// CachitaCache implements Cache with `gadelkareem/cachita` package.
type CachitaCache struct {
	store      cachita.Cache
	albumTag   string
	defaultTTL time.Duration
}

// NewCachitaCache returns a Cache service implemented using `gadelkareem/cachita`.
func NewCachitaCache() *CachitaCache {
	return &CachitaCache{
		store:      cachita.NewMemoryCache(1*time.Minute, 1*time.Minute),
		albumTag:   "album",
		defaultTTL: 60 * time.Minute,
	}
}

// Get reads an object data from the cache.
func (c *CachitaCache) GetAlbum(ctx context.Context, title string) (albums.Album, error) {
	item := albums.Album{}
	err := c.store.Get(c.albumKey(title), &item)
	if err == cachita.ErrNotFound {
		return albums.Album{}, ErrCacheMiss
	}
	return item, err
}

// Put store an object data to the cache.
func (c *CachitaCache) PutAlbum(ctx context.Context, album albums.Album) error {
	if err := c.store.Put(c.albumKey(album.Title), album, c.defaultTTL); err != nil {
		return err
	}
	return c.store.Tag(c.albumKey(album.Title), c.albumTag)
}

// InvalidateAlbum removes the specified Album from the cache.
func (c *CachitaCache) InvalidateAlbum(ctx context.Context, title string) error {
	return c.store.Invalidate(c.albumKey(title))
}

// InvalidateAllAlbums removes all albums from the cache.
func (c *CachitaCache) InvalidateAllAlbums(ctx context.Context) error {
	return c.store.InvalidateTags(c.albumTag)
}

// albumKey returns the cache key for an Album title.
func (c *CachitaCache) albumKey(title string) string {
	return c.albumTag + title
}
