package uploader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

// ResumableUploader implements resumable uploads.
// It will require a SessionStore to keep already started uploads.
type ResumableUploader struct {
	// HTTP Client
	client *http.Client
	// URL of the endpoint to upload to
	url string
	// store keeps active upload sessions.
	store SessionStore

	log log.Logger
}

// SessionStore represents an storage to keep resumable uploads.
type SessionStore interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}

// NewUploader returns an Uploader using the specified client or error in case
// of non valid configuration.
// The client must have the proper permissions to upload files.
//
// Use WithLogger(...) and WithEndpoint(...) to
// customize configuration.
func NewResumableUploader(client *http.Client, store SessionStore, options ...Option) (*ResumableUploader, error) {
	logger := defaultLogger()
	endpoint := defaultEndpoint()

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		case optkeyEndpoint:
			endpoint = o.Value().(string)
		}
	}

	u := &ResumableUploader{
		client: client,
		url:    endpoint,
		store:  store,
		log:    logger,
	}

	if store == nil {
		return nil, errors.New("store can't be nil")
	}

	return u, nil
}

// Upload returns the Google Photos upload token for an Upload object.
func (u *ResumableUploader) Upload(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating file upload: type=resumable, file=%s", upload.name)
	upload.sent = u.offsetFromPreviousSession(ctx, upload)

	if upload.sent == 0 {
		u.log.Debugf("Initiating new upload session: file=%s", upload.name)
		return u.createUploadSession(ctx, upload)
	}

	u.log.Debugf("Resuming previous upload session: file=%s", upload.name)
	return u.resumeUploadSession(ctx, upload)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *ResumableUploader) offsetFromPreviousSession(ctx context.Context, upload *Upload) int64 {
	// Get any previous session for this Upload
	url := string(u.store.Get(upload.fingerprint()))

	// Query previous upload status and get offset if active.
	req, err := upload.createQueryOffsetRequest(url)
	if err != nil {
		return 0
	}

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	status := res.Header.Get("X-Goog-Upload-Status")
	if status != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		u.store.Delete(upload.fingerprint())
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < upload.size {
		return offset
	}
	return 0
}

func (u *ResumableUploader) createUploadSession(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating upload session: file=%s", upload.name)
	req, err := upload.createInitialUploadRequest(u.url)
	if err != nil {
		return "", err
	}

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read upload url
	url := res.Header.Get("X-Goog-Upload-URL")
	u.store.Set(upload.fingerprint(), []byte(url))

	// Start upload session
	return u.resumeUploadSession(ctx, upload)
}

func (u *ResumableUploader) resumeUploadSession(ctx context.Context, upload *Upload) (string, error) {
	// Get any previous session for this Upload
	url := string(u.store.Get(upload.fingerprint()))

	_, err := upload.r.Seek(upload.sent, io.SeekStart)
	if err != nil {
		return "", err
	}

	req, err := upload.createResumeUploadRequest(url)
	if err != nil {
		return "", err
	}

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		u.log.Errorf("Failed to process request: err=%s", err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		u.log.Errorf("Failed to read response %s", err)
		return "", err
	}
	token := string(b)
	return token, nil
}

// createInitialUploadRequest returns a starting resumable upload request for Google Photos.
// url is the Google Photos upload endpoint.
func (u *Upload) createInitialUploadRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", path.Base(u.name))
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", u.size))

	return req, nil
}

// createQueryOffsetRequest returns a query offset request for Google Photos.
// url is the unique URL that must be used to complete the upload through all of the remaining requests.
func (u *Upload) createQueryOffsetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	return req, nil
}

// createResumeUploadRequest returns a resumable upload request to continue an started upload for Google Photos.
// url is the unique URL that must be used to complete the upload through all of the remaining requests.
func (u *Upload) createResumeUploadRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, u.r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Length", fmt.Sprintf("%d", u.size-u.sent))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", u.sent))

	return req, nil
}

// fingerprint returns a value to be used to identify upload session.
func (u *Upload) fingerprint() string {
	return fmt.Sprintf("%s|%d", u.name, u.size)
}