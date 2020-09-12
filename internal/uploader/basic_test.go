package uploader_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/api/googleapi"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/mock"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/uploader"
)

func TestNewBasicUploader(t *testing.T) {
	c := http.DefaultClient

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := uploader.NewBasicUploader(c)
		if err != nil {
			t.Errorf("should not return any error: err=%s", err)
		}
	})

	t.Run("WithLogger", func(t *testing.T) {
		want := &log.DiscardLogger{}
		_, err := uploader.NewBasicUploader(c, uploader.WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
	})

	t.Run("WithEndpoint", func(t *testing.T) {
		want := "https://my-domain.com/v1/upload"
		_, err := uploader.NewBasicUploader(c, uploader.WithEndpoint(want))
		if err != nil {
			t.Errorf("error was not expected at this point. err: %s", err)
		}
	})

	t.Run("WithEmptyEndpoint", func(t *testing.T) {
		_, err := uploader.NewBasicUploader(c, uploader.WithEndpoint(""))
		if err == nil {
			t.Errorf("error was expected")
		}
	})
}

func TestBasicUploader_Upload(t *testing.T) {
	i := mock.MockedUploadItem{}

	t.Run("ReturnsTokenOnSuccessfulUpload", func(t *testing.T) {
		want := "token"
		c := &mock.HttpClient{
			DoFn: func(req *http.Request) (response *http.Response, err error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader(want)),
				}, nil
			},
		}

		u, err := uploader.NewBasicUploader(c)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		got, err := u.Upload(context.Background(), i)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		if got != uploader.UploadToken(want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("ReturnsErrorOnFailedUpload", func(t *testing.T) {
		c := &mock.HttpClient{
			DoFn: func(req *http.Request) (response *http.Response, err error) {
				return nil, &googleapi.Error{Code: http.StatusTooManyRequests}
			},
		}

		u, err := uploader.NewBasicUploader(c)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		if _, err := u.Upload(context.Background(), i); err == nil {
			t.Errorf("error was expected.")
		}
	})
}
