package uploader

import (
	"bytes"
	"context"
	"io"
	"strings"
)

// MockedUploader mocks an uploading service.
type MockedUploader struct {
	UploadFileFn func(filepath string, ctx context.Context) (string, error)
	UploadFn     func(context.Context, UploadItem) (UploadToken, error)
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u MockedUploader) UploadFile(filepath string, ctx context.Context) (string, error) {
	return u.UploadFileFn(filepath, ctx)
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u MockedUploader) Upload(ctx context.Context, item UploadItem) (UploadToken, error) {
	return u.UploadFn(ctx, item)
}

// MockedUploadItem represents a mocked file upload item.
type MockedUploadItem struct {
	Path string
	size int64
}

// Open returns a io.ReadSeeker with a fixed string: "some test content inside a mocked file".
func (m MockedUploadItem) Open() (io.ReadSeeker, int64, error) {
	var b bytes.Buffer
	var err error

	r := strings.NewReader("some test content inside a mocked file")
	m.size, err = b.ReadFrom(r)
	if err != nil {
		return r, 0, err
	}
	return r, m.size, nil
}

// Name returns the name (path) of the item.
func (m MockedUploadItem) Name() string {
	return m.Path
}

// Size returns the length of "some test content inside a mocked file".
func (m MockedUploadItem) Size() int64 {
	return m.size
}

// MockedSessionStorer mocks a service to store resumable upload data.
type MockedSessionStorer struct {
	GetFn    func(f string) []byte
	SetFn    func(f string, u []byte)
	DeleteFn func(f string)
}

// Get invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Get(f string) []byte {
	return s.GetFn(f)
}

// Set invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Set(f string, u []byte) {
	s.SetFn(f, u)
}

// Delete invokes the mock implementation and marks the function as invoked.
func (s MockedSessionStorer) Delete(f string) {
	s.DeleteFn(f)
}
