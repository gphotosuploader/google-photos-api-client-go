package uploader

import (
	"context"
	"google.golang.org/api/googleapi"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
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

// NewSimpleUploader returns a new client to upload files to Google Photos.
// API methods require authentication, provide an [net/http.Client]
// that will perform the authentication for you (such as that provided
// by the [golang.org/x/oauth2] library).
func NewSimpleUploader(httpClient HttpClient) (*SimpleUploader, error) {
	defaultLogger := &log.DiscardLogger{}

	u := &SimpleUploader{
		client:  httpClient,
		BaseURL: defaultEndpoint,
		Logger:  defaultLogger,
	}

	return u, nil
}

// UploadFile uploads a file to Google Photos using upload request.
// A successful upload request returns an upload token. Use this upload
// token to create a media item with [media_items.Create].
func (u *SimpleUploader) UploadFile(ctx context.Context, filePath string) (uploadToken string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	upload, err := NewUploadFromFile(f)
	if err != nil {
		return "", err
	}

	return u.upload(ctx, upload)
}

func (u *SimpleUploader) upload(ctx context.Context, upload *Upload) (uploadToken string, err error) {
	req, err := u.createUploadRequest(upload)
	if err != nil {
		return "", err
	}

	res, err := u.doRequest(ctx, req)
	if err != nil {
		u.Logger.Errorf("Error while uploading %s: %s", upload, err)
		return "", err
	}
	defer res.Body.Close()

	return u.readUploadResponse(res, upload)
}

func (u *SimpleUploader) createUploadRequest(upload *Upload) (*http.Request, error) {
	req, err := http.NewRequest("POST", u.BaseURL, upload.stream)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", strconv.FormatInt(upload.size, 10))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", upload.Name)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	u.Logger.Debugf("Uploading %s (%d kB)", upload.Name, upload.size/1024)

	return req, nil
}

func (u *SimpleUploader) readUploadResponse(res *http.Response, upload *Upload) (uploadToken string, err error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		u.Logger.Errorf("Error while uploading %s: %s: could not read body: %s", upload, res.Status, err)
		return "", err
	}

	return string(body), nil
}

// doRequest executes the request call.
//
// Exactly one of *httpResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *httpResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header.
func (u *SimpleUploader) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}
