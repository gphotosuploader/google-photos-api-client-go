package gphotos

import (
	"io"
	"os"
	"path"
)

// FileUploadItem represents a local file.
type FileUploadItem string

// Open returns a stream.
// Caller should close it finally.
func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	f, err := os.Stat(m.String())
	if err != nil {
		return nil, 0, err
	}
	r, err := os.Open(m.String())
	if err != nil {
		return nil, 0, err
	}
	return r, f.Size(), nil
}

// Name returns the filename.
func (m FileUploadItem) Name() string {
	return path.Base(m.String())
}

func (m FileUploadItem) String() string {
	return string(m)
}

// Size returns size of the file.
func (m FileUploadItem) Size() int64 {
	f, err := os.Stat(m.String())
	if err != nil {
		return 0
	}
	return f.Size()
}
