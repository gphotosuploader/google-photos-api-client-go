package uploader

import (
	"fmt"
	"io"
	"os"
	"path"
)

// FileUploadItem represents a local file.
type FileUploadItem string

func NewFileUploadItem(filePath string) (FileUploadItem, error) {
	if !fileExists(filePath) {
		return "", fmt.Errorf("file does not exist (or is a directory")
	}
	return FileUploadItem(filePath), nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Open returns a stream.
// Caller should close it finally.
func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	r, err := os.Open(m.path())
	if err != nil {
		return nil, 0, err
	}
	fi, err := r.Stat()
	if err != nil {
		return nil, 0, err
	}
	return r, fi.Size(), nil
}

// Name returns the filename.
func (m FileUploadItem) Name() string {
	f, err := os.Stat(m.path())
	if err != nil {
		return ""
	}
	return path.Base(f.Name())
}

// Size returns size of the file.
func (m FileUploadItem) Size() int64 {
	f, err := os.Stat(m.path())
	if err != nil {
		return 0
	}
	return f.Size()
}

func (m FileUploadItem) path() string {
	return string(m)
}
