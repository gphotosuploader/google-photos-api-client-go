package cache

import (
	"errors"

	"github.com/gadelkareem/cachita"
)

// ErrCacheMiss is returned when a object is not found in cache.
var ErrCacheMiss = errors.New("gphotos: cache miss")

type Cache interface {
	albumsCache
}

// CachitaCache implements Cache using `gadelkareem/cachita` package.
type CachitaCache struct {
	cache cachita.Cache
}

func NewCachitaCache() *CachitaCache {
	return &CachitaCache{cache: cachita.Memory()}
}
