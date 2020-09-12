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

func WithPhotoService(s photoservice.Service) Option {
	return &option{
		name:  optkeyPhotoService,
		value: s,
	}
}

func WithUploader(u uploader.Uploader) Option {
	return &option{
		name:  optkeyUploader,
		value: u,
	}
}

func WithLogger(l log.Logger) Option {
	return &option{
		name:  optkeyLogger,
		value: l,
	}
}

func WithSessionStorer(s uploader.SessionStorer) Option {
	return &option{
		name:  optkeySessionStorer,
		value: s,
	}
}

func WithCacher(c cache.Cache) Option {
	return &option{
		name:  optkeyCacher,
		value: c,
	}
}
