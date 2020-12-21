package albums

import (
	"context"
	"errors"
	"time"

	"github.com/gadelkareem/cachita"
)

var (
	ErrCacheMiss = errors.New("item could not be found in the cache")
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

// GetAlbum reads an object data from the cache.
func (c CachitaCache) GetAlbum(ctx context.Context, title string) (Album, error) {
	item := Album{}
	err := c.store.Get(c.albumKey(title), &item)
	if err == cachita.ErrNotFound {
		return Album{}, ErrCacheMiss
	}
	return item, err
}

// PutAlbum store an object data to the cache.
func (c CachitaCache) PutAlbum(ctx context.Context, album Album) error {
	if err := c.store.Put(c.albumKey(album.Title), album, c.defaultTTL); err != nil {
		return err
	}
	return c.store.Tag(c.albumKey(album.Title), c.albumTag)
}

// PutManyAlbums store many objects data to the cache.
func (c CachitaCache) PutManyAlbums(ctx context.Context, albums []Album) error {
	for _, album := range albums {
		if err := c.PutAlbum(ctx, album); err != nil {
			return err
		}
	}
	return nil
}

// InvalidateAlbum removes the specified Album from the cache.
func (c CachitaCache) InvalidateAlbum(ctx context.Context, title string) error {
	return c.store.Invalidate(c.albumKey(title))
}

// InvalidateAllAlbums removes all albums from the cache.
func (c CachitaCache) InvalidateAllAlbums(ctx context.Context) error {
	return c.store.InvalidateTags(c.albumTag)
}

// albumKey returns the cache key for an Album title.
func (c CachitaCache) albumKey(title string) string {
	return c.albumTag + title
}
