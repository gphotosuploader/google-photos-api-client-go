package uploader

import (
	"io/ioutil"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi"
)

type mockedHttpClient struct {
	code int
	res *http.Response
}

func (c *mockedHttpClient) Do(req *http.Request) (*http.Response, error) {
	if c.code == http.StatusOK {
		return c.res, nil
	}
	err := new(googleapi.Error)
	err.Code = c.code
	return c.res, err
}

func responseAs(code int, header string, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Header:     nil,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
}
