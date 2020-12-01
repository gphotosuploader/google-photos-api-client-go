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

// Cache is used to store and retrieve previously obtained objects.
type Cache interface {
	// GetAlbum returns Album data from the cache corresponding to the specified title.
	// It will return ErrCacheMiss if there is no cached Album.
	GetAlbum(ctx context.Context, title string) (Album, error)

	// PutAlbum stores the Album data in the cache using the title as key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, GetAlbum, results in the original data.
	PutAlbum(ctx context.Context, album Album) error

	// PutManyAlbums stores many Album data in the cache using the title as key.
	PutManyAlbums(ctx context.Context, albums []Album) error

	// InvalidateAlbum removes the Album data from the cache corresponding to the specified title.
	// If there's no such Album in the cache, it will return nil.
	InvalidateAlbum(ctx context.Context, title string) error

	// InvalidateAllAlbums removes all key corresponding to albums
	InvalidateAllAlbums(ctx context.Context) error
}

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
