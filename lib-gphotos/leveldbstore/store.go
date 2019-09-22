package leveldbstore

import (
	"github.com/syndtr/goleveldb/leveldb"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
)

type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore create a new Store implemented by LevelDB
func NewLevelDBStore(path string) (gphotos.Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	store := &LevelDBStore{db: db}
	return store, err
}

// Get returns the url corresponding to the given fingerprint or error if not found
func (s *LevelDBStore) Get(fingerprint string) (string, error) {
	url, err := s.db.Get([]byte(fingerprint), nil)
	if err != nil {
		return "", err
	}
	return string(url), nil
}

// Set stores the url for a given fingerprint
func (s *LevelDBStore) Set(fingerprint, url string) error {
	return s.db.Put([]byte(fingerprint), []byte(url), nil)
}

func (s *LevelDBStore) Delete(fingerprint string) error {
	return s.db.Delete([]byte(fingerprint), nil)
}

func (s *LevelDBStore) Close() {
	s.Close()
}
