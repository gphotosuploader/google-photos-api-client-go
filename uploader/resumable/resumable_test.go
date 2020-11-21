package resumable

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewResumableUploader(t *testing.T) {
	c := &MockedHttpClient{}
	s := &MockedSessionStorer{}

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewResumableUploader(c, s)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		_, err := NewResumableUploader(c, s, WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := NewResumableUploader(c, s, WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithNilSessionStore", func(t *testing.T) {
		_, err := NewResumableUploader(c, nil)
		if err == nil {
			t.Errorf("error was expected when store in nil")
		}
	})
}
