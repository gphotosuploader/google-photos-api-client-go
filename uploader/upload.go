package uploader

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type Upload struct {
	stream io.ReadSeeker
	size   int64

	Name        string
	Fingerprint string
}

// NewUpload creates a new upload from an io.Reader.
func NewUpload(reader io.Reader, size int64, name string, fingerprint string) *Upload {
	stream, ok := reader.(io.ReadSeeker)

	if !ok {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(reader)
		if err != nil {
			return nil
		}
		stream = bytes.NewReader(buf.Bytes())
	}
	return &Upload{
		stream: stream,
		size:   size,

		Name:        name,
		Fingerprint: fingerprint,
	}
}

// NewUploadFromFile creates a new Upload from an os.File.
func NewUploadFromFile(f *os.File) (*Upload, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fingerprint := fmt.Sprintf("%s-%d-%s", fi.Name(), fi.Size(), fi.ModTime())

	return NewUpload(f, fi.Size(), fi.Name(), fingerprint), nil
}
