package gphotos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"google.golang.org/api/googleapi"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
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
	Get(fingerprint string) string
	Set(fingerprint string, url string)
	Delete(fingerprint string)
	Close()
}

// NewResumableUploader returns a new client to upload data to Google Photos
// with resumable capabilities.
// API methods require authentication, provide an [net/http.Client]
// that will perform the authentication for you (such as that provided
// by the [golang.org/x/oauth2] library).
func NewResumableUploader(httpClient HttpClient) (*ResumableUploader, error) {
	logger := &log.DiscardLogger{}
	endpoint := DefaultEndpoint

	u := &ResumableUploader{
		client:  httpClient,
		BaseURL: endpoint,
		Logger:  logger,
	}

	return u, nil
}

// UploadFile returns the Google Photos upload token after uploading a file.
// Any non-2xx status code is an error. Response headers are in error.(*googleapi.Error).Header.
func (u *ResumableUploader) UploadFile(ctx context.Context, filePath string) (UploadToken, error) {
	item, err := uploader.NewFileUploadItem(filePath)
	if err != nil {
		return "", err
	}

	u.Logger.Debugf("New resumable upload for file [%s].", item.Name())

	token, err := u.upload(ctx, item)
	return token, err
}

func (u *ResumableUploader) upload(ctx context.Context, item UploadItem) (UploadToken, error) {
	offset := u.offsetFromPreviousSession(ctx, item)
	u.Logger.Debugf("Current offset for [%s] is %d.", item.Name(), offset)
	if offset == 0 {
		return u.createUploadSession(ctx, item)
	}
	return u.resumeUploadSession(ctx, item, offset)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *ResumableUploader) offsetFromPreviousSession(ctx context.Context, item UploadItem) int64 {
	// Query previous upload status and get offsetFromResponse if active.
	if u.uploadSessionUrl(item) == "" {
		return 0
	}
	req, err := http.NewRequest("POST", u.uploadSessionUrl(item), nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")
	res, err := u.doRequest(ctx, req)
	if err != nil {
		return 0
	}
	defer res.Body.Close()
	return u.offsetFromResponse(res, item)
}

// offsetFromResponse returns the current offsetFromResponse if exist on the HTTP Response.
func (u *ResumableUploader) offsetFromResponse(res *http.Response, item UploadItem) int64 {
	if res.Header.Get("X-Goog-Upload-Status") != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		if u.isResumeEnabled() {
			u.Store.Delete(fingerprint(item))
		}
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < item.Size() {
		return offset
	}
	if u.isResumeEnabled() {
		u.Store.Delete(fingerprint(item))
	}
	return 0
}

func (u *ResumableUploader) createUploadSession(ctx context.Context, item UploadItem) (UploadToken, error) {
	req, err := u.prepareUploadRequest(item)
	if err != nil {
		return "", fmt.Errorf("creating upload session: %w", err)
	}

	res, err := u.doRequest(ctx, req)
	if err != nil {
		return "", fmt.Errorf("creating upload session: %w", err)
	}
	defer res.Body.Close()

	u.storeUploadSession(res, item)

	// Start upload session
	return u.resumeUploadSession(ctx, item, 0)
}

func (u *ResumableUploader) isResumeEnabled() bool {
	if u.Store != nil {
		return true
	}
	return false
}

// storeUploadSession keeps the upload session to allow resumes later.
func (u *ResumableUploader) storeUploadSession(res *http.Response, item UploadItem) {
	if url := res.Header.Get("X-Goog-Upload-URL"); url != "" && u.isResumeEnabled() {
		u.Store.Set(fingerprint(item), url)
	}
}

// prepareUploadRequest returns an HTTP request to upload item.
func (u *ResumableUploader) prepareUploadRequest(item UploadItem) (*http.Request, error) {
	_, size, err := item.Open()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.BaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", item.Name())
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", size))

	return req, nil
}

func (u *ResumableUploader) resumeUploadSession(ctx context.Context, item UploadItem, offset int64) (UploadToken, error) {
	u.Logger.Debugf("Resuming upload session for [%s] starting at offset %d", item.Name(), offset)
	req, err := u.prepareResumeUploadRequest(item, offset)
	if err != nil {
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	res, err := u.doRequest(ctx, req)
	if err != nil {
		u.Logger.Errorf("Failed to resume session: err=%s", err)
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		u.Logger.Errorf("Failed to read response %s", err)
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	token := string(b)
	return UploadToken(token), nil
}

func (u *ResumableUploader) prepareResumeUploadRequest(item UploadItem, offset int64) (*http.Request, error) {
	r, size, err := item.Open()
	if err != nil {
		return nil, fmt.Errorf("preparing resume upload request: %w", err)
	}
	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("preparing resume upload request: %w", err)
	}
	req, err := http.NewRequest("POST", u.uploadSessionUrl(item), r)
	if err != nil {
		return nil, fmt.Errorf("preparing resume upload request: %w", err)
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size-offset))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", offset))

	return req, nil
}

// doRequest executes the request call.
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

func (u *ResumableUploader) uploadSessionUrl(item UploadItem) string {
	if u.isResumeEnabled() {
		return u.Store.Get(fingerprint(item))
	}
	return ""
}

func fingerprint(item UploadItem) string {
	return fmt.Sprintf("%s|%d", item.Name(), item.Size())
}
