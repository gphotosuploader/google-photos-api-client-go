package gphotos

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

const (
	optkeyLogger        = "logger"
	optKeyCacher        = "cacher"
	optKeySessionStorer = "storer"
)

type Option interface {
	Name() string
	Value() interface{}
}

type option struct {
	name  string
	value interface{}
}

func (o option) Name() string       { return o.name }
func (o option) Value() interface{} { return o.value }

// WithLogger changes Client.log value.
func WithLogger(l log.Logger) Option {
	return &option{
		name:  optkeyLogger,
		value: l,
	}
}

func defaultLogger() log.Logger {
	return &log.DiscardLogger{}
}

// WithCacher changes Client.cache value.
func WithCacher(c cache.Cache) Option {
	return &option{
		name:  optKeyCacher,
		value: c,
	}
}

func defaultCacher() cache.Cache {
	return cache.NewCachitaCache()
}
