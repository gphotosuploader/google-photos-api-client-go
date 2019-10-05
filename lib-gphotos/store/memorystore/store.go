// Package memorystore provides implementation of in memory key/value database.
//
// Create or open a database:
//
// The returned DB instance is safe for concurrent use. Which mean that all
// DB's methods may be called concurrently from multiple goroutine.
//	db, err := memorystore.NewStore()
//	...
//	defer db.Close()
//	...
//
// Read or modify the database content:
//
//	// Remember that the contents of the returned slice should not be modified.
//	data := db.Get(key)
//	...
//	db.Put(key), []byte("value"))
//	...
//	db.Delete(key)
//	...
package memorystore

import (
	"sync"
)

// MemoryStore implements an in-memory Store.
type MemoryStore struct {
	m  map[string][]byte
	mu sync.RWMutex
}

// NewMemoryStore creates a new MemoryStore.
func NewStore() *MemoryStore {
	return &MemoryStore{
		m: make(map[string][]byte),
	}
}

func (s *MemoryStore) Get(key string) []byte {
	s.mu.RLock()
	value := s.m[key]
	s.mu.RUnlock()
	return value
}

func (s *MemoryStore) Set(key string, value []byte) {
	s.mu.Lock()
	s.m[key] = value
	s.mu.Unlock()
}

func (s *MemoryStore) Delete(key string) {
	s.mu.Lock()
	delete(s.m, key)
	s.mu.Unlock()
}

func (s *MemoryStore) Close() {
	for k := range s.m {
		delete(s.m, k)
	}
}
