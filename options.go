package gphotos

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
)

const (
	optkeyUploader      = "uploader"
	optKeyAlbumsService = "albumsRepo"
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

// WithUploader configures an Uploader.
func WithUploader(u uploader.MediaUploader) Option {
	return &option{
		name:  optkeyUploader,
		value: u,
	}
}

// WithUploader configures an Uploader.
func WithAlbums(u albums.AlbumsService) Option {
	return &option{
		name:  optKeyAlbumsService,
		value: u,
	}
}
