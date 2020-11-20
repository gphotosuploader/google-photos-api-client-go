package gphotos

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewClient(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected at this point: %s", err)
		}
	})
}
