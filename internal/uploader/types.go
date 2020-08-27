package uploader

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNilStore = errors.New("store can't be nil if Resume is enable")
)

type uploadService interface {
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

// SessionStore represents an storage to keep resumable uploads.
type SessionStore interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}
