package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

type mockUploadSessionStore struct{}

func (m *mockUploadSessionStore) Get(f string) []byte {
	return []byte(f)
}

func (m *mockUploadSessionStore) Set(f string, u []byte) {}

func (m *mockUploadSessionStore) Delete(f string) {}

func TestNewClientWithResumableUploads(t *testing.T) {
	c := http.DefaultClient
	store := &mockUploadSessionStore{}

	t.Run("EmptyHTTPClient", func(t *testing.T) {
		_, err := gphotos.NewClientWithResumableUploads(nil, store)
		if err == nil {
			t.Errorf("NewClientWithResumableUploads error was expected here")
		}
	})

	t.Run("WithNilUploadSessionStore", func(t *testing.T) {
		_, err := gphotos.NewClientWithResumableUploads(c, nil)
		if err != uploader.ErrNilStore {
			t.Errorf("NewClientWithResumableUploads - error was expected here: got=%s, want=%s", err, uploader.ErrNilStore)
		}
	})

	t.Run("WithoutOptions", func(t *testing.T) {
		got, err := gphotos.NewClientWithResumableUploads(c, store)
		if err != nil {
			t.Errorf("NewClientWithResumableUploads - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithResumableUploads - Photos service was not created")
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		l := &log.DiscardLogger{}
		got, err := gphotos.NewClientWithResumableUploads(c, store, gphotos.WithLogger(l))
		if err != nil {
			t.Errorf("NewClientWithResumableUploads - error was not expected here: err=%s", err)
		}
		if got.Service == nil {
			t.Errorf("NewClientWithResumableUploads - Photos service was not created")
		}
	})
}