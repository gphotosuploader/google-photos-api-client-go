package gphotos_test

import (
	"errors"
	"testing"

	"google.golang.org/api/googleapi"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/gphotos"
)

func TestIsRetryableError(t *testing.T) {

	t.Run("WithRetryableErrors", func(t *testing.T) {
		for code := 500; code <= 599; code++ {
			err := error(&googleapi.Error{Code: code})

			if !gphotos.IsRetryableError(err) {
				t.Errorf("error %d should be retryable.", code)
			}
		}
	})

	t.Run("WithNonRetryableErrors", func(t *testing.T) {
		for code := 400; code <= 499; code++ {
			err := error(&googleapi.Error{Code: code})

			if gphotos.IsRetryableError(err) {
				t.Errorf("error %d should not be retryable.", code)
			}
		}
	})

	t.Run("WithNonGoogleApiError", func(t *testing.T) {
		err := errors.New("an error")
		if !gphotos.IsRetryableError(err) {
			t.Errorf("a non Google API error should be retryable.")
		}
	})
}

func TestIsRateLimitError(t *testing.T) {
	t.Run("WithErrorDueToRateLimit", func(t *testing.T) {
		err := error(&googleapi.Error{Code:429})

		if !gphotos.IsRateLimitError(err) {
			t.Errorf("error 429 is due to rate limit.")
		}
	})

	t.Run("WithErrorNotDueToRateLimit", func(t *testing.T) {
		err := error(&googleapi.Error{Code:404})

		if gphotos.IsRateLimitError(err) {
			t.Errorf("error 404 is not due to rate limit.")
		}
	})

	t.Run("WithNonGoogleApiError", func(t *testing.T) {
		err := errors.New("an error")
		if gphotos.IsRateLimitError(err) {
			t.Errorf("a non Google API error should not be treat as rate limit caused.")
		}
	})
}
