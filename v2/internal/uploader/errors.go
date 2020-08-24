package uploader

import "errors"

var (
	ErrNilStore = errors.New("store can't be nil if Resume is enable")
)
