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
	Uploader Uploader
	r        io.ReadSeeker
	name     string
	size     int64
	sent     int64
	atEOF    bool
}

func (u *Upload) getUploadToken(ctx context.Context) (string, error) {
	if u.Uploader.resume {
		return u.getUploadTokenResumable(ctx)
	}
	return u.getUploadTokenNonResumable(ctx)
}

// getUploadTokenResumable uploads media data and returns the Google Photos token.
func (u *Upload) getUploadTokenResumable(ctx context.Context) (string, error) {
	u.Uploader.log.Printf("[DEBUG] Initiating file upload: type=resumable, file=%s", u.name)
	u.sent = u.offsetFromPreviousSession(ctx)

	if u.sent == 0 {
		u.Uploader.log.Printf("[DEBUG] Initiating new upload session: file=%s", u.name)
		return u.createUploadSession(ctx)
	}

	u.Uploader.log.Printf("[DEBUG] Resuming previous upload session: file=%s", u.name)
	return u.resumeUploadSession(ctx)
}

// getUploadToken upload media content and returns the Google Photos UploadToken.
func (u *Upload) getUploadTokenNonResumable(ctx context.Context) (string, error) {
	u.Uploader.log.Printf("[DEBUG] Initiating file upload: type=non-resumable, file=%s", u.name)

	req, err := http.NewRequest("POST", u.Uploader.url, u.r)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("X-Goog-Upload-File-Name", path.Base(u.name))
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	res, err := u.Uploader.c.Do(req.WithContext(ctx))
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

func (u *Upload) fingerprint() string {
	return fmt.Sprintf("%s|%d", u.name, u.size)
}

func (u *Upload) offsetFromPreviousSession(ctx context.Context) int64 {
	// Get any previous session for this Upload
	url, _ := u.Uploader.store.Get(u.fingerprint())

	// Query previous upload status and get offset if active.
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "query")

	res, err := u.Uploader.c.Do(req.WithContext(ctx))
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	status := res.Header.Get("X-Goog-Upload-Status")
	if status != "active" {
		// Other known statuses "final" and "cancelled" are both considered as already completed.
		// Let's restart the upload from scratch.
		_ = u.Uploader.store.Delete(u.fingerprint())
		return 0
	}

	offset, err := strconv.ParseInt(res.Header.Get("X-Goog-Upload-Size-Received"), 10, 64)
	if err == nil && offset > 0 && offset < u.size {
		return offset
	}
	return 0
}

func (u *Upload) resumeUploadSession(ctx context.Context) (string, error) {
	// Get any previous session for this Upload
	url, _ := u.Uploader.store.Get(u.fingerprint())

	_, err := u.r.Seek(u.sent, io.SeekStart)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, u)
	if err != nil {
		u.Uploader.log.Printf("[ERR] Failed to prepare request: err=%s", err)
		return "", err
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", u.size-u.sent))
	req.Header.Add("X-Goog-Upload-Command", "upload, finalize")
	req.Header.Set("X-Goog-Upload-Offset", fmt.Sprintf("%d", u.sent))

	res, err := u.Uploader.c.Do(req.WithContext(ctx))
	if err != nil {
		u.Uploader.log.Printf("[ERR] Failed to process request: err=%s", err)
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		u.Uploader.log.Printf("[ERR] Failed to read response %s", err)
		return "", err
	}
	token := string(b)
	return token, nil
}

func (u *Upload) createUploadSession(ctx context.Context) (string, error) {
	u.Uploader.log.Printf("[DEBUG] Initiating upload session: file=%s", u.name)
	req, err := http.NewRequest("POST", u.Uploader.url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Goog-Upload-Command", "start")
	req.Header.Add("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Protocol", "resumable")
	req.Header.Set("X-Goog-Upload-Raw-Size", fmt.Sprintf("%d", u.size))

	res, err := u.Uploader.c.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Read upload url
	url := res.Header.Get("X-Goog-Upload-URL")
	_ = u.Uploader.store.Set(u.fingerprint(), url)

	// Start upload session
	return u.resumeUploadSession(ctx)
}
