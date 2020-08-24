package gphotos

import (
	"context"
	"errors"
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// ErrCacheMiss is returned when a object is not found in cache.
var ErrCacheMiss = errors.New("photos: cache miss")

// Cache is used to store and retrieve previously obtained objects.
type Cache interface {
	// GetAlbum returns Album data for the specified key.
	// If there's no such key, GetAlbum returns ErrCacheMiss.
	GetAlbum(ctx context.Context, key string) (*photoslibrary.Album, error)

	// PutAlbum stores the Album data in the cache under the specified key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, GetAlbum, results in the original data.
	PutAlbum(ctx context.Context, key string, album *photoslibrary.Album, ttl time.Duration) error

	// DeleteAlbum removes the Album data from the cache under the specified key.
	// If there's no such key in the cache, DeleteAlbum returns nil.
	InvalidateAlbum(ctx context.Context, key string) error
}

// CachitaCache implements Cache using `gadelkareem/cachita` package.
type CachitaCache struct {
	cache cachita.Cache
}

func NewCachitaCache() *CachitaCache {
	return &CachitaCache{cache: cachita.Memory()}
}

// Get reads an object data from the cache.
func (m *CachitaCache) GetAlbum(ctx context.Context, key string) (*photoslibrary.Album, error) {
	i := &photoslibrary.Album{}
	err := m.cache.Get(prefixAlbumKey(key), i)
	if err == cachita.ErrNotFound {
		return nil, ErrCacheMiss
	}

	return i, err
}

// Put store an object data to the cache.
func (m *CachitaCache) PutAlbum(ctx context.Context, key string, album *photoslibrary.Album, ttl time.Duration) error {
	return m.cache.Put(prefixAlbumKey(key), *album, ttl)
}

// InvalidateAlbum removes the specified Album from the cache.
func (m *CachitaCache) InvalidateAlbum(ctx context.Context, key string) error {
	return m.cache.Invalidate(prefixAlbumKey(key))
}

func prefixAlbumKey(key string) string {
	return "album:" + key
}
