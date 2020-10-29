package gphotos

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/cache"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/photoservice"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

const (
	optkeyLogger        = "logger"
	optkeyCacher        = "cacher"
	optkeySessionStorer = "sessionStorer"
	optkeyPhotoService  = "photoservice"
	optkeyUploader      = "uploader"
)

// Option represents a configurable parameter for Google Photos API client.
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

// WithPhotoService configures a Google Photos service.
func WithPhotoService(s photoservice.Service) Option {
	return &option{
		name:  optkeyPhotoService,
		value: s,
	}
}

// WithUploader configures an Uploader.
func WithUploader(u uploader.Uploader) Option {
	return &option{
		name:  optkeyUploader,
		value: u,
	}
}

// WithLogger configures a Logger.
func WithLogger(l log.Logger) Option {
	return &option{
		name:  optkeyLogger,
		value: l,
	}
}

// WithSessionStorer configures a service to keep resumable uploads.
func WithSessionStorer(s uploader.SessionStorer) Option {
	return &option{
		name:  optkeySessionStorer,
		value: s,
	}
}

// WithCacher configures a Cache.
func WithCacher(c cache.Cache) Option {
	return &option{
		name:  optkeyCacher,
		value: c,
	}
}
