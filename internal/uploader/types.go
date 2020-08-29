package uploader

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNilStore = errors.New("store can't be nil if Resume is enable")
)

type Uploader interface {
	// Upload uploads the media item. It returns an upload token.
	Upload(context.Context, UploadItem) (UploadToken, error)
}

// UploadItem represents an uploadable item.
type UploadItem interface {
	// Open returns a stream.
	// Caller should close it finally.
	Open() (io.ReadCloser, int64, error)
	// Name returns the filename.
	Name() string
	// String returns the full name, e.g. path or URL.
	String() string
}

// UploadToken represents a pointer to the uploaded item.
type UploadToken string

// Upload represents an object to be uploaded.
type Upload struct {
	r    io.ReadSeeker
	name string
	size int64
	sent int64
}

