package uploader

import (
	"errors"
)

var (
	ErrUploadNotFound    = errors.New("upload not found")
	ErrFingerprintNotSet = errors.New("fingerprint not set")
)
