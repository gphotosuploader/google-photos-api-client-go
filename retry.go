package gphotos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically, so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the BaseURL is invalid. This error isn't typed
	// specifically, so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by Google Photos when
	// the request quota limit has been exceeded. This error isn't typed
	// specifically, so we resort to matching on the error string.
	requestQuotaErrorRe = regexp.MustCompile(`Quota exceeded for quota metric 'All requests' and limit 'All requests per day'`)
)

// addRetryHandler returns an HTTP client with a retry policy.
func addRetryHandler(client *http.Client) *http.Client {
	c := retryablehttp.Client{
		HTTPClient: client,

		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     3,

		CheckRetry: GooglePhotosServiceRetryPolicy,

		Backoff: retryablehttp.DefaultBackoff,
	}

	return c.StandardClient()
}

// GooglePhotosServiceRetryPolicy provides a retry policy implementing Google Photos
// best practices.
//
// See: https://developers.google.com/photos/library/guides/best-practices#error-handling
func GooglePhotosServiceRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	shouldRetry, err := baseRetryPolicy(resp, err)

	var e *ErrDailyQuotaExceeded
	if errors.As(err, &e) {
		return false, err
	}

	// don't propagate other errors
	return shouldRetry, nil
}

func baseRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(urlErr.Error()) {
				return false, urlErr
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(urlErr.Error()) {
				return false, urlErr
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// 429 Too Many Requests can be recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from a client.
	// If the 'write requests per minute per user' quota is exceeded, the error is recoverable.
	// If the 'daily API' quota is exceeded, the error is not recoverable.
	if resp.StatusCode == http.StatusTooManyRequests {
		slurp, ioerr := io.ReadAll(resp.Body)
		if ioerr != nil {
			return false, ioerr
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(slurp))

		// Don't retry if the 'All request' per day quota has been exceeded.
		if requestQuotaErrorRe.MatchString(string(slurp)) {
			return false, &ErrDailyQuotaExceeded{}
		}

		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	return false, nil
}
