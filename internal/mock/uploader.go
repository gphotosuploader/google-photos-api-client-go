package mock

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader"
)

// Uploader mocks an uploading service.
type Uploader struct {
	UploadFileFn func(filepath string, ctx context.Context) (string, error)
	UploadFileInvoked bool

	UploadFn      func(context.Context, uploader.UploadItem) (uploader.UploadToken, error)
	UploadInvoked bool
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u *Uploader) UploadFile(filepath string, ctx context.Context) (string, error) {
	u.UploadFileInvoked = true
	return u.UploadFileFn(filepath, ctx)
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u *Uploader) Upload(ctx context.Context, item uploader.UploadItem) (uploader.UploadToken, error) {
	u.UploadInvoked = true
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

// SessionStorer mocks a service to store resumable upload data.
type SessionStorer struct {
	GetFn      func(f string) []byte
	GetInvoked bool

	SetFn      func(f string, u []byte)
	SetInvoked bool

	DeleteFn      func(f string)
	DeleteInvoked bool
}

// Get invokes the mock implementation and marks the function as invoked.
func (s *SessionStorer) Get(f string) []byte {
	s.GetInvoked = true
	return s.GetFn(f)
}

// Set invokes the mock implementation and marks the function as invoked.
func (s *SessionStorer) Set(f string, u []byte) {
	s.SetInvoked = true
	s.SetFn(f, u)
}

// Delete invokes the mock implementation and marks the function as invoked.
func (s *SessionStorer) Delete(f string) {
	s.DeleteInvoked = true
	s.DeleteFn(f)
}
