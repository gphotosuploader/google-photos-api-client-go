package photoservice

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

const (
	optkeyLogger = "logger"
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

// WithLogger changes GooglePhotosService.log value.
func WithLogger(l log.Logger) Option {
	return &option{
		name:  optkeyLogger,
		value: l,
	}
}

func defaultLogger() log.Logger {
	return &log.DiscardLogger{}
}