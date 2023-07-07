package gphotos_test

import (
	"context"
	"errors"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGooglePhotoServiceRetryPolicy(t *testing.T) {
	testCases := []struct {
		name            string
		body            string
		statusCode      int
		shouldBeRetried bool
		expectedError   error
	}{
		// SHOULD BE RETRIED
		{
			name:            "TooManyRequest response should retry (except 'Daily requests per day exceeded')",
			body:            ` `,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: true,
			expectedError:   nil,
		},
		{
			name:            "TooManyRequest for 'Write requests per minute exceeded' response should retry",
			body:            sampleGoogleWriteRequestsPerMinuteExceededBodyResponse,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: true,
			expectedError:   nil,
		},
		{
			name:            "InternalServerError response should retry",
			body:            ` `,
			statusCode:      http.StatusInternalServerError,
			shouldBeRetried: true,
			expectedError:   nil,
		},
		// SHOULD NOT BE RETRIED
		{
			name:            "Ok response should not retry",
			body:            ` `,
			statusCode:      http.StatusOK,
			shouldBeRetried: false,
			expectedError:   nil,
		},
		{
			name:            "TooManyRequest for 'Daily requests per day exceeded' response should not retry",
			body:            mocks.SampleGoogleRequestPerDayExceededBodyResponse,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: false,
			expectedError:   &gphotos.ErrDailyQuotaExceeded{},
		},
		{
			name:            "BadRequest response should not retry",
			body:            ` `,
			statusCode:      http.StatusBadRequest,
			shouldBeRetried: false,
			expectedError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: tc.statusCode,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}
			got, err := gphotos.GooglePhotosServiceRetryPolicy(context.Background(), res, nil)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error")
			}

			if tc.shouldBeRetried != got {
				t.Errorf("want: %t, got: %t", tc.shouldBeRetried, got)
			}
		})
	}
}

func TestContextErr_Retry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(5 * time.Millisecond)

	// ctx was already cancelled.
	shouldRetry, err := gphotos.GooglePhotosServiceRetryPolicy(ctx, &http.Response{}, nil)
	if err == nil {
		t.Fatalf("error was expected at this point but not happened")
	}

	if shouldRetry {
		t.Errorf("should not retry")
	}
}

func TestURLErrors_Retry(t *testing.T) {
	testCases := []struct {
		name            string
		err             error
		shouldBeRetried bool
	}{
		{
			name:            "Too many redirects should not retry",
			err:             &url.Error{Err: errors.New("stopped after 10 redirects")},
			shouldBeRetried: false,
		},
		{
			name:            "Invalid URL schema should not retry",
			err:             &url.Error{Err: errors.New("unsupported protocol scheme")},
			shouldBeRetried: false,
		},
		{
			name:            "Any other error should be retried",
			err:             &url.Error{Err: errors.New("a different error")},
			shouldBeRetried: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gphotos.GooglePhotosServiceRetryPolicy(context.Background(), &http.Response{}, tc.err)
			if err != nil {
				t.Fatalf("error was not expected at this point: %v", err)
			}

			if tc.shouldBeRetried != got {
				t.Errorf("want: %t, got: %t", tc.shouldBeRetried, got)
			}
		})
	}
}

const sampleGoogleWriteRequestsPerMinuteExceededBodyResponse = `
{
  "error": {
    "code": 429,
    "message": "Quota exceeded for quota metric 'Write requests' and limit 'Write requests per minute per user' of service 'photoslibrary.googleapis.com' for consumer 'project_number:844831818923'.",
    "errors": [
      {
        "message": "Quota exceeded for quota metric 'Write requests' and limit 'Write requests per minute per user' of service 'photoslibrary.googleapis.com' for consumer 'project_number:844831818923'.",
        "domain": "global",
        "reason": "rateLimitExceeded"
      }
    ],
    "status": "RESOURCE_EXHAUSTED",
    "details": [
      {
        "@type": "type.googleapis.com/google.rpc.ErrorInfo",
        "reason": "RATE_LIMIT_EXCEEDED",
        "domain": "googleapis.com",
        "metadata": {
          "service": "photoslibrary.googleapis.com",
          "quota_limit_value": "30",
          "quota_location": "global",
          "consumer": "projects/844831818923",
          "quota_metric": "photoslibrary.googleapis.com/write_requests",
          "quota_limit": "WritesPerMinutePerUser"
        }
      },
      {
        "@type": "type.googleapis.com/google.rpc.Help",
        "links": [
          {
            "description": "Request a higher quota limit.",
            "url": "https://cloud.google.com/docs/quota#requesting_higher_quota"
          }
        ]
      }
    ]
  }
}
`
