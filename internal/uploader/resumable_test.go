package uploader

import (
	"net/http"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewUploader(t *testing.T) {
	c := http.DefaultClient

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewUploader(c)
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		l := &log.DiscardLogger{}
		_, err := NewUploader(c, WithLogger(l))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := NewUploader(c, WithEndpointURL(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionResumableUploads", func(t *testing.T) {
		store := &mockUploadSessionStore{}
		got, err := NewUploader(c, WithResumableUploads(store))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}

		if !got.resume {
			t.Errorf("NewUploader resumable uploads were expected here")
		}
	})

	t.Run("WithOptionResumableUploadsNilStore", func(t *testing.T) {
		_, err := NewUploader(c, WithResumableUploads(nil))
		if err == nil {
			t.Errorf("NewUploader error was expected here")
		}
	})
}

type mockUploadSessionStore struct{}

func (m *mockUploadSessionStore) Get(f string) []byte {
	return []byte(f)
}

func (m *mockUploadSessionStore) Set(f string, u []byte) {}

func (m *mockUploadSessionStore) Delete(f string) {}
