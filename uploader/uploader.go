package uploader

import (
	"net/http"
)

// defaultEndpoint is the Google Photos endpoint for uploads.
const defaultEndpoint = "https://photoslibrary.googleapis.com/v1/uploads"

// HttpClient represent a client to make an HTTP request.
// It is usually implemented by [net/http.Client].
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// UploadToken represents a pointer to the uploaded item in Google Photos.
// Use this upload token to create a media item with [media_items.Create].
type UploadToken string
