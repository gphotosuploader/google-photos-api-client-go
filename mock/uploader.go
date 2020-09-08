package mock

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

type Uploader struct {
	UploadFn      func(context.Context, uploader.UploadItem) (uploader.UploadToken, error)
	UploadInvoked bool
}

// Upload invokes the mock implementation and marks the function as invoked.
func (u *Uploader) Upload(ctx context.Context, item uploader.UploadItem) (uploader.UploadToken, error) {
	u.UploadInvoked = true
	return u.UploadFn(ctx, item)
}

// FileUploadItem represents a mocked file upload item.
type FileUploadItem struct {
	Path string
	size int64
}

func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	var b bytes.Buffer
	var err error

	r := strings.NewReader("some test content inside a mocked file")
	m.size, err = b.ReadFrom(r)
	if err != nil {
		return r, 0, err
	}
	return r, m.size, nil
}

func (m FileUploadItem) Name() string {
	return m.String()
}

func (m FileUploadItem) String() string {
	return m.Path
}

func (m FileUploadItem) Size() int64 {
	return m.size
}

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
