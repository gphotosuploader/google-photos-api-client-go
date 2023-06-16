package albums

import "errors"

var (
	// ErrAlbumNotFound is the error returned when an album is not found.
	ErrAlbumNotFound = errors.New("album not found")
)
