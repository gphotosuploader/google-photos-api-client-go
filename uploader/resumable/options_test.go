package resumable

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v3/internal/log"
)

func TestWithLogger(t *testing.T) {
	want := &log.DiscardLogger{}

	got := WithLogger(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestWithEndpoint(t *testing.T) {
	want := "https://domain.com/uploads"

	got := WithEndpoint(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
