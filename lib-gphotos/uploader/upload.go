package uploader

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
)

// Upload represents an object to be uploaded.
type Upload struct {
	r        io.ReadSeeker
	name     string
	size     int64
	sent     int64
	finished bool
}

func NewUpload(r io.ReadSeeker, name string, size int64) *Upload {
	return &Upload{
		r:    r,
		name: name,
		size: size,
	}
}

// Upload returns the Google Photos upload token for an Upload object.
func (u *Uploader) Upload(ctx context.Context, upload *Upload) (string, error) {
	if u.resume {
		return u.uploadWithResumeCapability(ctx, upload)
	}
	return u.uploadWithoutResumeCapability(ctx, upload)
}

// uploadWithResumeCapability returns the Google Photos upload token using resume uploads to upload data.
func (u *Uploader) uploadWithResumeCapability(ctx context.Context, upload *Upload) (string, error) {
	u.log.Printf("[DEBUG] Initiating file upload: type=resumable, file=%s", upload.name)
	upload.sent = u.offsetFromPreviousSession(ctx, upload)

	if upload.sent == 0 {
		u.log.Printf("[DEBUG] Initiating new upload session: file=%s", upload.name)
		return u.createUploadSession(ctx, upload)
	}

	u.log.Printf("[DEBUG] Resuming previous upload session: file=%s", upload.name)
	return u.resumeUploadSession(ctx, upload)
}

// upload returns the Google Photos upload token using non-resumable upload.
func (u *Uploader) uploadWithoutResumeCapability(ctx context.Context, upload *Upload) (string, error) {
	u.log.Printf("[DEBUG] Initiating file upload: type=non-resumable, file=%s", upload.name)

	r := &ReadProgressReporter{
		r:        upload.r,
		filename: upload.name,
		size:     upload.size,
		sent:     upload.sent,
		logger:   u.log,
	}
	req, err := http.NewRequest("POST", u.url, r)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("X-Goog-Upload-File-Name", path.Base(upload.name))
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	res, err := u.c.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	token := string(b)
	return token, nil
}

// fingerprint returns a value to be used to identify upload session.
func (u *Upload) fingerprint() string {
	return fmt.Sprintf("%s|%d", u.name, u.size)
}

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *Uploader) offsetFromPreviousSession(ctx context.Context, upload *Upload) int64 {
	// Get any previous session for this Upload
	url, _ := u.store.Get(upload.fingerprint())

	// Query previous upload status and get offset if active.
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	res, err := u.c.Do(req.WithContext(ctx))
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	status := res.Header.Get("X-Goog-Upload-Status")
	if status != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		_ = u.store.Delete(upload.fingerprint())
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < upload.size {
		return offset
	}
	return 0
}

func (u *Uploader) resumeUploadSession(ctx context.Context, upload *Upload) (string, error) {
	// Get any previous session for this Upload
	url, _ := u.store.Get(upload.fingerprint())

	_, err := upload.r.Seek(upload.sent, io.SeekStart)
	if err != nil {
		return "", err
	}

	r := &ReadProgressReporter{
		r:        upload.r,
		filename: upload.name,
		size:     upload.size,
		sent:     upload.sent,
		logger:   u.log,
	}
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		u.log.Printf("[ERR] Failed to prepare request: err=%s", err)
		return "", err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", upload.size-upload.sent))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", upload.sent))

	res, err := u.c.Do(req.WithContext(ctx))
	if err != nil {
		u.log.Printf("[ERR] Failed to process request: err=%s", err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		u.log.Printf("[ERR] Failed to read response %s", err)
		return "", err
	}
	token := string(b)
	return token, nil
}

func (u *Uploader) createUploadSession(ctx context.Context, upload *Upload) (string, error) {
	u.log.Printf("[DEBUG] Initiating upload session: file=%s", upload.name)
	req, err := http.NewRequest("POST", u.url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Add("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", upload.size))

	res, err := u.c.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read upload url
	url := res.Header.Get("X-Goog-Upload-URL")
	_ = u.store.Set(upload.fingerprint(), url)

	// Start upload session
	return u.resumeUploadSession(ctx, upload)
}
