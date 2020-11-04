package uploader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

// ResumableUploader implements resumable uploads.
// It will require a SessionStorer to keep already upload session.
type ResumableUploader struct {
	// client is an HTTP client used for uploading. It needs the proper authentication in place.
	client httpClient
	// url is the url the endpoint to upload to
	url string
	// store is an upload session store.
	store SessionStorer
	// log is a logger to send messages.
	log log.Logger
}

// SessionStorer represents an storage to keep resumable uploads.
type SessionStorer interface {
	Get(fingerprint string) []byte
	Set(fingerprint string, url []byte)
	Delete(fingerprint string)
}

// NewResumableUploader returns an Uploader or error in case of non valid configuration.
// The supplied client must have the proper authentication to upload files.
// The supplied store will be used to keep upload sessions.
//
// Use WithLogger(...) and WithEndpoint(...) to customize configuration.
func NewResumableUploader(client httpClient, store SessionStorer, options ...Option) (*ResumableUploader, error) {
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

	// validate configuration options.
	if store == nil {
		return nil, errors.New("session store can't be nil")
	}

	return u, nil
}

// UploadFile returns the Google Photos upload token after uploading a file.
func (u *ResumableUploader) UploadFile(filePath string, ctx context.Context) (string, error)  {
	token, err := u.Upload(ctx, FileUploadItem(filePath))
	return string(token), err
}

// Upload returns the Google Photos upload token for an Upload object.
func (u *ResumableUploader) Upload(ctx context.Context, item UploadItem) (UploadToken, error) {
	u.log.Debugf("Initiating file upload: type=resumable, file=%s", item.Name())
	offset := u.offsetFromPreviousSession(ctx, item)

	if offset == 0 {
		u.log.Debugf("Initiating new upload session: file=%s", item.Name())
		return u.createUploadSession(ctx, item)
	}

	u.log.Debugf("Resuming previous upload session: file=%s", item.Name())
	return u.resumeUploadSession(ctx, item, offset)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *ResumableUploader) offsetFromPreviousSession(ctx context.Context, item UploadItem) int64 {
	// Query previous upload status and get offsetFromResponse if active.
	req, err := http.NewRequest("POST", u.uploadSessionUrl(item), nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	return u.offsetFromResponse(res, item)
}

// offsetFromResponse returns the current offsetFromResponse if exist on the HTTP Response.
func (u *ResumableUploader) offsetFromResponse(res *http.Response, item UploadItem) int64 {
	status := res.Header.Get("X-Goog-Upload-Status")
	if status != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		u.store.Delete(fingerprint(item))
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < item.Size() {
		return offset
	}
	return 0
}

func (u *ResumableUploader) createUploadSession(ctx context.Context, item UploadItem) (UploadToken, error) {
	u.log.Debugf("Initiating upload session: file=%s", item.Name())

	req, err := u.prepareUploadRequest(item)
	if err != nil {
		return "", err
	}

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	u.storeUploadSession(res, item)

	// Start upload session
	return u.resumeUploadSession(ctx, item, 0)
}

// storeUploadSession keeps the upload session to allow resumes later.
func (u *ResumableUploader) storeUploadSession(res *http.Response, item UploadItem) {
	if url := res.Header.Get("X-Goog-Upload-URL"); url != "" {
		u.store.Set(fingerprint(item), []byte(url))
	}
}

// prepareUploadRequest returns an HTTP request to upload item.
func (u *ResumableUploader) prepareUploadRequest(item UploadItem) (*http.Request, error) {
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

func (u *ResumableUploader) resumeUploadSession(ctx context.Context, item UploadItem, offset int64) (UploadToken, error) {
	req, err := u.prepareResumeUploadRequest(item, offset)
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
	return UploadToken(token), nil
}

func (u *ResumableUploader) prepareResumeUploadRequest(item UploadItem, offset int64) (*http.Request, error) {
	r, size, err := item.Open()
	if err != nil {
		return nil, err
	}

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.uploadSessionUrl(item), r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size-offset))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", offset))

	return req, nil
}

func (u *ResumableUploader) uploadSessionUrl(item UploadItem) string {
	return string(u.store.Get(fingerprint(item)))
}

func fingerprint(item UploadItem) string {
	return fmt.Sprintf("%s|%d", item.Name(), item.Size())
}
