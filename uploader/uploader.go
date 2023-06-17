package gphotos

import (
	"io"
	"net/http"
)

// DefaultEndpoint is the Google Photos endpoint for uploads.
const DefaultEndpoint = "https://photoslibrary.googleapis.com/v1/uploads"

// HttpClient represent a client to make an HTTP request.
// It is usually implemented by [net/http.Client].
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// UploadToken represents a pointer to the uploaded item in Google Photos.
// Use this upload token to create a media item with [media_items.Create].
type UploadToken string

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
