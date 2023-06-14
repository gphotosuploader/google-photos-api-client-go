package gphotos

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

var (
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically, so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically, so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

	// A regular expression to match the error returned by Google Photos when
	// the request quota limit has been exceeded. This error isn't typed
	// specifically, so we resort to matching on the error string.
	requestQuotaErrorRe = regexp.MustCompile(`Quota exceeded for quota metric 'All requests' and limit 'All requests per day'`)

	// A regular expression to match the error returned by Google Photos when
	// the storage quota limit has been exceeded. This error isn't typed
	// specifically, so we resort to matching on the error string.
	storageQuotaErrorRe = regexp.MustCompile(`The remaining storage in the user's account is not enough to perform this operation`)
)

// newRetryHandler returns an HTTP client with a retry policy.
func newRetryHandler(client *http.Client) *http.Client {
	c := retryablehttp.NewClient()
	c.HTTPClient = client
	c.CheckRetry = defaultRetryPolicy
	c.Logger = nil // Disable DEBUG logs
	return c.StandardClient()
}

// defaultRetryPolicy provides a default callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
func defaultRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// don't propagate other errors
	shouldRetry, _ := baseRetryPolicy(resp, err)
	return shouldRetry, nil
}

func baseRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to a requests quota limit exceed.
			if requestQuotaErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to a storage quota limit exceed.
			if storageQuotaErrorRe.MatchString(v.Error()) {
				return false, v
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, v
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// 429 Too Many Requests can be recoverable. Sometimes the server puts
	// a Retry-After response header to indicate when the server is
	// available to start processing request from a client.
	// If the write requests per minute per user quota is exceeded, the error is recoverable.
	// If the daily API quota is exceeded, the error is not recoverable.
	if resp.StatusCode == http.StatusTooManyRequests {
		slurp, ioerr := io.ReadAll(resp.Body)
		if ioerr != nil {
			return false, ioerr
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(slurp))

		if requestQuotaErrorRe.MatchString(string(slurp)) {
			return false, fmt.Errorf("daily quota exceeded")
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
