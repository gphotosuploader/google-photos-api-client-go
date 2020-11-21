package basic

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewBasicUploader(t *testing.T) {
	c := &MockedHttpClient{}

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewBasicUploader(c)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		_, err := NewBasicUploader(c, WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := NewBasicUploader(c, WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})
}

