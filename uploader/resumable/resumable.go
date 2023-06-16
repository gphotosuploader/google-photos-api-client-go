package resumable

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"google.golang.org/api/googleapi"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
)

// ResumableUploader implements resumable uploads.
// It will require a SessionStorer to keep already upload session.
type ResumableUploader struct {
	// authenticatedClient is an HTTP client used for uploading. It needs the proper authentication in place.
	authenticatedClient HttpClient
	// url is the url the endpoint to upload to
	url string
	// store is an upload session store.
	store SessionStorer
	// log is a logger to send messages.
	log log.Logger
}

// HttpClient represents a HTTP client.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SessionStorer represents an storage to keep resumable uploads.
type SessionStorer interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}

// NewResumableUploader returns an Uploader or error in case of non valid configuration.
// The supplied authenticatedClient must have the proper authentication to upload files.
// The supplied store will be used to keep upload sessions.
//
// Use WithLogger(...) and WithEndpoint(...) to customize configuration.
func NewResumableUploader(authenticatedClient HttpClient, store SessionStorer, options ...Option) (*ResumableUploader, error) {
	logger := defaultLogger()
	endpoint := uploader.DefaultEndpoint

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optkeyEndpoint:
			endpoint = o.Value().(string)
		}
	}

	u := &ResumableUploader{
		authenticatedClient: authenticatedClient,
		url:                 endpoint,
		store:               store,
		log:                 logger,
	}

	// validate configuration options.
	if store == nil {
		return nil, errors.New("session store can't be nil")
	}

	return u, nil
}

// UploadFile returns the Google Photos upload token after uploading a file.
// Any non-2xx status code is an error. Response headers are in error.(*googleapi.Error).Header.
func (u ResumableUploader) UploadFile(ctx context.Context, filePath string) (string, error) {
	item, err := uploader.NewFileUploadItem(filePath)
	if err != nil {
		return "", err
	}
	u.log.Debugf("New resumable upload for file [%s].", item.Name())
	token, err := u.upload(ctx, item)
	return string(token), err
}

func (u ResumableUploader) upload(ctx context.Context, item uploader.UploadItem) (uploader.UploadToken, error) {
	offset := u.offsetFromPreviousSession(ctx, item)
	u.log.Debugf("Current offset for [%s] is %d.", item.Name(), offset)
	if offset == 0 {
		return u.createUploadSession(ctx, item)
	}
	return u.resumeUploadSession(ctx, item, offset)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u ResumableUploader) offsetFromPreviousSession(ctx context.Context, item uploader.UploadItem) int64 {
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
func (u ResumableUploader) offsetFromResponse(res *http.Response, item uploader.UploadItem) int64 {
	if res.Header.Get("X-Goog-Upload-Status") != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		u.store.Delete(fingerprint(item))
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < item.Size() {
		return offset
	}
	u.store.Delete(fingerprint(item))
	return 0
}

func (u ResumableUploader) createUploadSession(ctx context.Context, item uploader.UploadItem) (uploader.UploadToken, error) {
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

// storeUploadSession keeps the upload session to allow resumes later.
func (u ResumableUploader) storeUploadSession(res *http.Response, item uploader.UploadItem) {
	if url := res.Header.Get("X-Goog-Upload-URL"); url != "" {
		u.store.Set(fingerprint(item), []byte(url))
	}
}

// prepareUploadRequest returns an HTTP request to upload item.
func (u ResumableUploader) prepareUploadRequest(item uploader.UploadItem) (*http.Request, error) {
	_, size, err := item.Open()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.url, nil)
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

func (u ResumableUploader) resumeUploadSession(ctx context.Context, item uploader.UploadItem, offset int64) (uploader.UploadToken, error) {
	u.log.Debugf("Resuming upload session for [%s] starting at offset %d", item.Name(), offset)
	req, err := u.prepareResumeUploadRequest(item, offset)
	if err != nil {
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	res, err := u.doRequest(ctx, req)
	if err != nil {
		u.log.Errorf("Failed to resume session: err=%s", err)
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		u.log.Errorf("Failed to read response %s", err)
		return "", fmt.Errorf("resuming upload session: %w", err)
	}
	token := string(b)
	return uploader.UploadToken(token), nil
}

func (u ResumableUploader) prepareResumeUploadRequest(item uploader.UploadItem, offset int64) (*http.Request, error) {
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
func (u ResumableUploader) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := u.authenticatedClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (u ResumableUploader) uploadSessionUrl(item uploader.UploadItem) string {
	return string(u.store.Get(fingerprint(item)))
}

func fingerprint(item uploader.UploadItem) string {
	return fmt.Sprintf("%s|%d", item.Name(), item.Size())
}
