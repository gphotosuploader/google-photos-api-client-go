package memorystore_test

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos/store/memorystore"
)

func TestSetWithInMemoryDB(t *testing.T) {
	var testData = []struct {
		key   string
		value string
		want  string
	}{
		{"", "", ""},
		{"", "testValue", "testValue"},
		{"testKey", "", ""},
		{"testKey", "testValue", "testValue"},
	}
	s := memorystore.NewStore()
	defer s.Close()

	for i, tt := range testData {
		s.Set(tt.key, []byte(tt.value))
		got := string(s.Get(tt.key))

		if got != tt.want {
			t.Errorf("test case failed: id=%d, got=%s, want=%s", i, got, tt.want)
		}
	}
}

func TestDeleteWithInMemory(t *testing.T) {
	s := memorystore.NewStore()
	defer s.Close()

	t.Run("DeleteNonExistentKey", func(t *testing.T) {
		s.Delete("testKey")
	})

	t.Run("DeleteExistentKey", func(t *testing.T) {
		s.Set("testKey", []byte("testValue"))
		s.Delete("testKey")
	})

	t.Run("DeleteEmptyKey", func(t *testing.T) {
		key, value := "testKey", "testValue"
		s.Set(key, []byte(value))

		// Delete an empty key, and see if it only affects one key
		s.Delete("")

		got := string(s.Get(key))
		if got != value {
			t.Errorf("delete case failed: got=%s, want=%s", got, value)
		}
	})
}

func TestGetWithInMemory(t *testing.T) {
	s := memorystore.NewStore()
	defer s.Close()

	t.Run("GetNonExistentKey", func(t *testing.T) {
		got := string(s.Get("non-existant"))
		if got != "" {
			t.Errorf("get case failed: got=%s, want=\"\"", got)
		}
	})

	t.Run("GetExistentKey", func(t *testing.T) {
		key, value := "testKey", "testValue"
		s.Set(key, []byte(value))
		got := string(s.Get(key))
		if got != value {
			t.Errorf("get case failed: got=%s, want=%s", got, value)
		}

	})
}
