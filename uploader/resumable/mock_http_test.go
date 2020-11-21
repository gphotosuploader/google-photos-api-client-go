package resumable

import (
	"net/http"
)

// HttpClient mocks an HTTP client.
type MockedHttpClient struct {
	DoFn func(req *http.Request) (*http.Response, error)
}

// Do invokes the mock implementation and marks the function as invoked.
func (c MockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoFn(req)
}
