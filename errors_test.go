package gphotos_test

import (
	"context"
	"errors"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"net/http"
	"testing"
)

func TestErrDailyQuotaExceeded_Error(t *testing.T) {
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	httpClient := http.DefaultClient
	client, err := gphotos.NewClientWithBaseURL(httpClient, srv.URL())
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}

	_, err = client.Albums.GetById(context.Background(), mocks.ShouldReachDailyQuota)

	var e *gphotos.ErrDailyQuotaExceeded
	if !errors.As(err, &e) {
		t.Errorf("unexpected error: %v", err)
	}
}
