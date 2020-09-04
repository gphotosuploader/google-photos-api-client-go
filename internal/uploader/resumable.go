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
		return nil, errors.New("store can't be nil")
	}

	return u, nil
}

// Upload returns the Google Photos upload token for an Upload object.
func (u *ResumableUploader) Upload(ctx context.Context, item UploadItem) (UploadToken, error) {
	u.log.Debugf("Initiating file upload: type=resumable, file=%s", item.String())
	offset := u.offsetFromPreviousSession(ctx, item)

	if offset == 0 {
		u.log.Debugf("Initiating new upload session: file=%s", item.String())
		return u.createUploadSession(ctx, item)
	}

	u.log.Debugf("Resuming previous upload session: file=%s", item.String())
	return u.resumeUploadSession(ctx, item, offset)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *ResumableUploader) offsetFromPreviousSession(ctx context.Context, item UploadItem) int64 {
	// Query previous upload status and get offset if active.
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
	u.log.Debugf("Initiating upload session: file=%s", item.String())

	_, size, err := item.Open()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", u.url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", item.Name())
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", size))

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read upload url
	if url := res.Header.Get("X-Goog-Upload-URL"); url != "" {
		u.store.Set(fingerprint(item), []byte(url))
	}

	// Start upload session
	return u.resumeUploadSession(ctx, item, 0)
}

func (u *ResumableUploader) resumeUploadSession(ctx context.Context, item UploadItem, offset int64) (UploadToken, error) {
	r, size, err := item.Open()
	if err != nil {
		return "", nil
	}

	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", u.uploadSessionUrl(item), r)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size-offset))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", offset))

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

func (u *ResumableUploader) uploadSessionUrl(item UploadItem) string {
	return string(u.store.Get(fingerprint(item)))
}

func fingerprint(item UploadItem) string {
	return fmt.Sprintf("%s|%d", item.Name(), item.Size())
}
