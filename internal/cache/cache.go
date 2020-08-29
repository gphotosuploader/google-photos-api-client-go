package cache

import (
	"context"
	"errors"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// ErrCacheMiss is returned when a object is not found in cache.
var ErrCacheMiss = errors.New("gphotos: cache miss")

type Cache interface {
	albumsCache
}

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
