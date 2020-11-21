package uploader

import (
	"context"
	"io"
)

const DefaultEndpoint = "https://photoslibrary.googleapis.com/v1/uploads"

type MediaUploader interface {
	UploadFile(ctx context.Context, filePath string) (uploadToken string, err error)
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




