package uploader

import (
	"bytes"
	"io"
	"strings"
)

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
