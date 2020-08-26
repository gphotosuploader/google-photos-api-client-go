package cache

import (
	"context"
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// albumsCache is used to store and retrieve previously obtained Albums.
type albumsCache interface {
	// GetAlbum returns Album data from the cache corresponding to the specified key.
	// It will return ErrCacheMiss if there is no cached Album with that key.
	GetAlbum(ctx context.Context, key string) (*photoslibrary.Album, error)

	// PutAlbum stores the Album data in the cache under the specified key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, GetAlbum, results in the original data.
	PutAlbum(ctx context.Context, key string, album *photoslibrary.Album, ttl time.Duration) error

	// DeleteAlbum removes the Album data from the cache under the specified key.
	// If there's no such key in the cache, DeleteAlbum returns nil.
	InvalidateAlbum(ctx context.Context, key string) error
}

// encodeAlbumKey returns the cache key for an Album.
func encodeAlbumKey(key string) string {
	return "album:" + key
}

// Get reads an object data from the cache.
func (m *CachitaCache) GetAlbum(ctx context.Context, key string) (*photoslibrary.Album, error) {
	i := &photoslibrary.Album{}
	err := m.cache.Get(encodeAlbumKey(key), i)
	if err == cachita.ErrNotFound {
		return nil, ErrCacheMiss
	}

	return i, err
}

// Put store an object data to the cache.
func (m *CachitaCache) PutAlbum(ctx context.Context, key string, album *photoslibrary.Album, ttl time.Duration) error {
	return m.cache.Put(encodeAlbumKey(key), *album, ttl)
}

// InvalidateAlbum removes the specified Album from the cache.
func (m *CachitaCache) InvalidateAlbum(ctx context.Context, key string) error {
	return m.cache.Invalidate(encodeAlbumKey(key))
}
