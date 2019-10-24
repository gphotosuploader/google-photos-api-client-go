package uploader

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

// UploadFromFile returns the Google Photos upload token for the file.
func (u *Uploader) UploadFromFile(ctx context.Context, filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed opening file: err=%s", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed getting file maxBytes: file=%s, err=%s", filename, err)
	}
	size := fileStat.Size()

	upload := NewUpload(file, filename, size)
	return u.Upload(ctx, upload)

}

// Upload returns the Google Photos upload token for an Upload object.
func (u *Uploader) Upload(ctx context.Context, upload *Upload) (string, error) {
	if u.resume {
		return u.uploadWithResumeCapability(ctx, upload)
	}
	return u.uploadWithoutResumeCapability(ctx, upload)
}

// Upload represents an object to be uploaded.
type Upload struct {
	r    io.ReadSeeker
	name string
	size int64
	sent int64
}

func NewUpload(r io.ReadSeeker, name string, size int64) *Upload {
	return &Upload{
		r:    r,
		name: name,
		size: size,
	}
}

// fingerprint returns a value to be used to identify upload session.
func (u *Upload) fingerprint() string {
	return fmt.Sprintf("%s|%d", u.name, u.size)
}

// uploadWithResumeCapability returns the Google Photos upload token using resume uploads to upload data.
func (u *Uploader) uploadWithResumeCapability(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating file upload: type=resumable, file=%s", upload.name)
	upload.sent = u.offsetFromPreviousSession(ctx, upload)

	if upload.sent == 0 {
		u.log.Debugf("Initiating new upload session: file=%s", upload.name)
		return u.createUploadSession(ctx, upload)
	}

	u.log.Debugf("Resuming previous upload session: file=%s", upload.name)
	return u.resumeUploadSession(ctx, upload)
}

// upload returns the Google Photos upload token using non-resumable upload.
func (u *Uploader) uploadWithoutResumeCapability(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating file upload: type=non-resumable, file=%s", upload.name)

	req, err := createRawUploadRequest(u.url, upload)
	if err != nil {
		return "", err
	}

	res, err := u.client.Do(req.WithContext(ctx))
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

// offsetFromPreviousSession returns the bytes already uploaded in previous upload sessions.
func (u *Uploader) offsetFromPreviousSession(ctx context.Context, upload *Upload) int64 {
	// Get any previous session for this Upload
	url := string(u.store.Get(upload.fingerprint()))

	// Query previous upload status and get offset if active.
	req, err := createQueryOffsetRequest(url)
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

func (u *Uploader) resumeUploadSession(ctx context.Context, upload *Upload) (string, error) {
	// Get any previous session for this Upload
	url := string(u.store.Get(upload.fingerprint()))

	_, err := upload.r.Seek(upload.sent, io.SeekStart)
	if err != nil {
		return "", err
	}

	req, err := createResumeUploadRequest(url, upload)
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

func (u *Uploader) createUploadSession(ctx context.Context, upload *Upload) (string, error) {
	u.log.Debugf("Initiating upload session: file=%s", upload.name)
	req, err := createInitialResumableUploadRequest(u.url, upload)
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
