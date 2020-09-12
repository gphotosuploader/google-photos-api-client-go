package uploader

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

const (
	optkeyLogger   = "logger"
	optkeyEndpoint = "endpoint"
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

// WithEndpoint changes the Client.endpoint value.
func WithEndpoint(u string) Option {
	return &option{
		name:  optkeyEndpoint,
		value: u,
	}
}

func defaultLogger() log.Logger {
	return &log.DiscardLogger{}
}

func defaultEndpoint() string {
	return "https://photoslibrary.googleapis.com/v1/uploads"
}
