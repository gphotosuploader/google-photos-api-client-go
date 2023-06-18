package gphotos_test

import (
	"context"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGooglePhotoServiceRetryPolicy(t *testing.T) {
	testCases := []struct {
		name            string
		body            string
		statusCode      int
		shouldBeRetried bool
	}{
		// SHOULD BE RETRIED
		{
			name:            "TooManyRequest response should retry (except 'Daily requests per day exceeded')",
			body:            ` `,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: true,
		},
		{
			name:            "TooManyRequest for 'Write requests per minute exceeded' response should retry",
			body:            sampleGoogleWriteRequestsPerMinuteExceededBodyResponse,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: true,
		},
		{
			name:            "InternalServerError response should retry",
			body:            ` `,
			statusCode:      http.StatusInternalServerError,
			shouldBeRetried: true,
		},

		// SHOULD NOT BE RETRIED
		{
			name:            "Ok response should not retry",
			body:            ` `,
			statusCode:      http.StatusOK,
			shouldBeRetried: false,
		},
		{
			name:            "TooManyRequest for 'Daily requests per day exceeded' response should not retry",
			body:            sampleGoogleRequestPerDayExceededBodyResponse,
			statusCode:      http.StatusTooManyRequests,
			shouldBeRetried: false,
		},
		{
			name:            "BadRequest response should not retry",
			body:            ` `,
			statusCode:      http.StatusBadRequest,
			shouldBeRetried: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: tc.statusCode,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}
			got, _ := gphotos.GooglePhotosServiceRetryPolicy(context.Background(), res, nil)

			if tc.shouldBeRetried != got {
				t.Errorf("want: %t, got: %t", tc.shouldBeRetried, got)
			}
		})
	}
}

const sampleGoogleRequestPerDayExceededBodyResponse = `
{
  "error": {
    "code": 429,
    "message": "Quota exceeded for quota metric 'All requests' and limit 'All requests per day' of service 'photoslibrary.googleapis.com' for consumer 'project_number:844831818923'.",
    "errors": [
      {
        "message": "Quota exceeded for quota metric 'All requests' and limit 'All requests per day' of service 'photoslibrary.googleapis.com' for consumer 'project_number:844831818923'.",
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
          "quota_limit_value": "10000",
          "consumer": "projects/844831818923",
          "service": "photoslibrary.googleapis.com",
          "quota_limit": "ApiCallsPerProjectPerDay",
          "quota_location": "global",
          "quota_metric": "photoslibrary.googleapis.com/all_requests"
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
