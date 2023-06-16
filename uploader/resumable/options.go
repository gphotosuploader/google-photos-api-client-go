package resumable

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
)

const (
	optkeyLogger   = "logger"
	optkeyEndpoint = "endpoint"
)

// Option represents a configurable parameter of the Uploader.
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
func WithLogger(l log.Logger) *option {
	return &option{
		name:  optkeyLogger,
		value: l,
	}
}

// WithEndpoint changes the Client.endpoint value.
func WithEndpoint(u string) *option {
	return &option{
		name:  optkeyEndpoint,
		value: u,
	}
}

func defaultLogger() log.Logger {
	return &log.DiscardLogger{}
}
