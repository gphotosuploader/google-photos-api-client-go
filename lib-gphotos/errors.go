package gphotos

import "golang.org/x/xerrors"

var (
	ErrNilStore = xerrors.New("store can't be nil if Resume is enable")
)
