package gphotos_test

import (
	"net/http"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
)

func TestNewClient(t *testing.T) {
	_, err := gphotos.NewClient(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

}
