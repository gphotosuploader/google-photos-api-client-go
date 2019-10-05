package store

// Store interface could be implemented by several key-value databases.
// In this package are included:
// - LevelDB implementation
// - In Memory implementation
//
// Once the database is created:
//
// Read or modify the database content:
//
// Remember that the contents of the returned slice should not be modified.
//	data := db.Get(key)
//	...
//	db.Put(key), []byte("value"))
//	...
//	db.Delete(key)
//	...
type Store interface {
	Get(key string) []byte
	Set(key string, value []byte)
	Delete(key string)
	Close()
}

