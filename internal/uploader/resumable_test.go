package uploader_test

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

func TestNewResumableUploader(t *testing.T) {
	c := &mock.HttpClient{}
	s := &mock.SessionStorer{}

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := uploader.NewResumableUploader(c, s)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		_, err := uploader.NewResumableUploader(c, s, uploader.WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := uploader.NewResumableUploader(c, s, uploader.WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithNilSessionStore", func(t *testing.T) {
		_, err := uploader.NewResumableUploader(c, nil)
		if err == nil {
			t.Errorf("error was expected when store in nil")
		}
	})
}
