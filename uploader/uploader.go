package uploader

import (
	"context"
	"io"
	"net/http"
)

type MediaUploader interface {
	UploadFile(filePath string, ctx context.Context) (string, error)
}

// UploadItem represents an uploadable item.
type UploadItem interface {
	// Open returns a stream.
	// Caller should close it finally.
	Open() (io.ReadSeeker, int64, error)
	// Name returns the filename.
	Name() string
	// Size returns the size (in bytes).
	Size() int64
}

// UploadToken represents a pointer to the uploaded item.
type UploadToken string

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
