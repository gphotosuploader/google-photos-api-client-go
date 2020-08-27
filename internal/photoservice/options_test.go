package photoservice_test

import (
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/gphotos"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestWithLogger(t *testing.T) {
	want := &log.DiscardLogger{}

	got := photoservice.WithLogger(want)
	if got.Value() != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
