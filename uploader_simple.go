package gphotos

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
)

// SimpleUploader implements a simple uploader to Google Photos.
type SimpleUploader struct {
	client HttpClient // HTTP client used to communicate with the API.

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// Logger used to log messages.
	Logger log.Logger
}

// HttpClient represent a client to make an HTTP request.
// It is usually implemented by [net/http.Client].
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// UploadToken represents a pointer to the uploaded item in Google Photos.
// Use this upload token to create a media item with [media_items.Create].
type UploadToken string

// NewSimpleUploader returns a new client to upload data to Google Photos.
// API methods require authentication, provide an [net/http.Client]
// that will perform the authentication for you (such as that provided
// by the [golang.org/x/oauth2] library).
func NewSimpleUploader(httpClient HttpClient) (*SimpleUploader, error) {
	var defaultLogger = &log.DiscardLogger{}

	u := &SimpleUploader{
		client:  httpClient,
		BaseURL: uploader.DefaultEndpoint,
		Logger:  defaultLogger,
	}

	return u, nil
}

// UploadFile upload bytes to Google Photos using upload requests.
// A successful upload request returns an upload token. Use this upload
// token to create a media item with [media_items.Create].
func (u *SimpleUploader) UploadFile(ctx context.Context, filePath string) (string, error) {
	token, err := u.upload(ctx, uploader.FileUploadItem(filePath))
	return string(token), err
}

func (u *SimpleUploader) upload(ctx context.Context, uploadItem uploader.UploadItem) (UploadToken, error) {
	req, err := u.prepareUploadRequest(uploadItem)
	if err != nil {
		return "", err
	}

	u.Logger.Debugf("Uploading %s (%d kB)", uploadItem.Name(), uploadItem.Size()/1024)

	res, err := u.client.Do(req)
	if err != nil {
		u.Logger.Errorf("Error while uploading %s: %s", uploadItem, err)
		return "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		u.Logger.Errorf("Error while uploading %s: %s: could not read body: %s", uploadItem, res.Status, err)
		return "", err
	}
	body := string(b)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %s: %s", res.Status, body)
	}

	return UploadToken(body), nil

}

// prepareUploadRequest returns an HTTP request to upload item.
//
// See: https://developers.google.com/photos/library/guides/upload-media#uploading-bytes.
func (u *SimpleUploader) prepareUploadRequest(item uploader.UploadItem) (*http.Request, error) {
	r, size, err := item.Open()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.BaseURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", item.Name())
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	return req, nil
}
