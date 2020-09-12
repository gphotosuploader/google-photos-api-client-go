package uploader_test

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

func TestWithLogger(t *testing.T) {
	want := &log.DiscardLogger{}

	got := uploader.WithLogger(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithEndpoint(t *testing.T) {
	want := "https://domain.com/uploads"

	got := uploader.WithEndpoint(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
