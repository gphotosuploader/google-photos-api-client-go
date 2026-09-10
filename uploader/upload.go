package uploader

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Upload struct {
	stream io.ReadSeeker
	size   int64

	Name        string
	Fingerprint string
	SHA256      string
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

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, fmt.Errorf("hashing upload: %w", err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("rewinding upload: %w", err)
	}
	contentHash := hex.EncodeToString(h.Sum(nil))
	fingerprintInput := filepath.Clean(f.Name()) + "\x00" + strconv.FormatInt(fi.Size(), 10) + "\x00" + contentHash
	fingerprintHash := sha256.Sum256([]byte(fingerprintInput))
	upload := NewUpload(f, fi.Size(), fi.Name(), hex.EncodeToString(fingerprintHash[:]))
	upload.SHA256 = contentHash
	return upload, nil
}
