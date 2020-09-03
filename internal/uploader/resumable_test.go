package uploader

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewResumableUploader(t *testing.T) {
	c := &mockedHttpClient{}
	s := &mockUploadSessionStore{}

	t.Run("WithoutOptions", func(t *testing.T) {
		got, err := NewResumableUploader(c, s)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
		if got.store != s {
			t.Errorf("want: %v, got: %v", s, got.log)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		got, err := NewResumableUploader(c, s, WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
		if got.log != want {
			t.Errorf("want: %v, got: %v", want, got.log)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		got, err := NewResumableUploader(c, s, WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
		if got.url != want {
			t.Errorf("want: %v, got: %v", want, got.url)
		}
	})

	t.Run("WithNilSessionStore", func(t *testing.T) {
		_, err := NewResumableUploader(c, nil)
		if err == nil {
			t.Errorf("error was expected when store in nil")
		}
	})
}

type mockUploadSessionStore struct{}

func (m *mockUploadSessionStore) Get(f string) []byte {
	return []byte(f)
}

func (m *mockUploadSessionStore) Set(f string, u []byte) {}

func (m *mockUploadSessionStore) Delete(f string) {}
