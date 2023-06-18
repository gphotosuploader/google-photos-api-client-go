package uploader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"google.golang.org/api/googleapi"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
)

// ResumableUploader implements resumable uploads.
//
// See: https://developers.google.com/photos/library/guides/resumable-uploads.
type ResumableUploader struct {
	client HttpClient // HTTP client used to communicate with the API.

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// Logger used to log messages.
	Logger log.Logger

	// Store maps an upload's fingerprint with the corresponding upload URL.
	Store Store
}

// Store represents a service to map upload's fingerprint with
// the corresponding upload URL.
type Store interface {
	Get(fingerprint string) (string, bool)
	Set(fingerprint string, url string)
	Delete(fingerprint string)
	Close()
}

// NewResumableUploader returns a new client to upload files to Google Photos
// with resumable capabilities.
// API methods require authentication, provide an [net/http.Client]
// that will perform the authentication for you (such as that provided
// by the [golang.org/x/oauth2] library).
func NewResumableUploader(httpClient HttpClient) (*ResumableUploader, error) {
	defaultLogger := &log.DiscardLogger{}

	u := &ResumableUploader{
		client:  httpClient,
		BaseURL: defaultEndpoint,
		Logger:  defaultLogger,
	}

	return u, nil
}

// UploadFile returns the Google Photos upload token after uploading a file.
// Any non-2xx status code is an error. Response headers are in error.(*googleapi.Error).Header.
func (u *ResumableUploader) UploadFile(ctx context.Context, filePath string) (uploadToken string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	upload, err := NewUploadFromFile(f)
	if err != nil {
		return "", err
	}

	u.Logger.Debugf("Starting resumable upload for file [%s].", filePath)

	return u.createOrResumeUpload(ctx, upload)
}

func (u *ResumableUploader) createOrResumeUpload(ctx context.Context, upload *Upload) (uploadToken string, err error) {
	uploadToken, err = u.resumeUpload(ctx, upload)

	if err == nil {
		return uploadToken, err
	}

	return u.createUpload(ctx, upload)
}

func (u *ResumableUploader) createUpload(ctx context.Context, upload *Upload) (uploadToken string, err error) {
	req, err := http.NewRequest("POST", u.BaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Name", upload.Name)
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", strconv.FormatInt(upload.size, 10))

	res, err := u.doRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("creating upload: %w", err)
	}
	defer res.Body.Close()

	if http.StatusOK == res.StatusCode {
		location := res.Header.Get("X-Goog-Upload-URL")

		if u.isResumeEnabled() {
			u.Store.Set(upload.Fingerprint, location)
		}
	}

	return u.resumeUpload(ctx, upload)
}

func (u *ResumableUploader) isResumeEnabled() bool {
	return u.Store != nil
}

func (u *ResumableUploader) resumeUpload(ctx context.Context, upload *Upload) (uploadToken string, err error) {
	if len(upload.Fingerprint) == 0 {
		return "", ErrFingerprintNotSet
	}
	url, found := u.Store.Get(upload.Fingerprint)

	if !found {
		return "", ErrUploadNotFound
	}

	offset, err := u.getUploadOffset(ctx, url)
	if err != nil {
		return "", err
	}

	if _, err := upload.stream.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, upload.stream)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Length", strconv.FormatInt(upload.size, 10))
	req.Header.Set("X-Goog-Upload-Offset", strconv.FormatInt(offset, 10))
	req.Header.Set("X-Goog-Upload-Command", "upload, finalize")

	res, err := u.doRequest(ctx, req)
	if err != nil {
		u.Logger.Errorf("Failed to resume upload: %s", err)
		return "", fmt.Errorf("resuming upload: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		u.Logger.Errorf("Failed to read response: %s", err)
		return "", fmt.Errorf("resuming upload: %w", err)
	}

	if u.isResumeEnabled() {
		u.Store.Delete(upload.Fingerprint)
	}

	return string(b), nil
}

func (u *ResumableUploader) getUploadOffset(ctx context.Context, url string) (int64, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	res, err := u.doRequest(ctx, req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	if res.Header.Get("X-Goog-Upload-Status") != "active" {
		// Other known statuses "final" and "canceled" are both considered as already completed.
		// Let's restart the upload from scratch.
		return -1, ErrUploadNotFound
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err != nil {
		return -1, fmt.Errorf("getUploadOffset: %w", err)

	}

	return offset, nil
}

// doRequest executes the request call.
//
// Exactly one of *httpResponse or error will be non-nil.
// Any non-2xx status code is an error. Response headers are in either
// *httpResponse.Header or (if a response was returned at all) in
// error.(*googleapi.Error).Header.
func (u *ResumableUploader) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (u *ResumableUploader) IsResumeEnabled() bool {
	return u.isResumeEnabled()

}
