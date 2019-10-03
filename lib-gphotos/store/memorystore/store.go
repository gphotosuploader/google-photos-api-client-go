package memorystore

import gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"

// MemoryStore implements an in-memory Store.
type MemoryStore struct {
	m map[string]string
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() (gphotos.Store, error) {
	return &MemoryStore{
		make(map[string]string),
	}, nil
}

func (s *MemoryStore) Get(fingerprint string) (string, bool) {
	url, ok := s.m[fingerprint]
	return url, ok
}

func (s *MemoryStore) Set(fingerprint, url string) {
	s.m[fingerprint] = url
}

func (s *MemoryStore) Delete(fingerprint string) {
	delete(s.m, fingerprint)
}

func (s *MemoryStore) Close() {
	for k := range s.m {
		delete(s.m, k)
	}
}
