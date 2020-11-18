package uploader

import (
	"net/http"
)

// HttpClient mocks an HTTP client.
type HttpClient struct {
	DoFn func(req *http.Request) (*http.Response, error)
}

// Do invokes the mock implementation and marks the function as invoked.
func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoFn(req)
}
