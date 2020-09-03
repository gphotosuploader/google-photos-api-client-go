package uploader

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewBasicUploader(t *testing.T) {
	c := &mockedHttpClient{}

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := NewBasicUploader(c)
		if err != nil {
			t.Errorf("should not return any error: err=%s", err)
		}
	})

	t.Run("WithLogger", func(t *testing.T) {
		want := &log.DiscardLogger{}
		got, err := NewBasicUploader(c, WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if got.log != want {
			t.Errorf("want: %v, got: %v", want, got.log)
		}
	})

	t.Run("WithEndpoint", func(t *testing.T) {
		want := "https://my-domain.com/v1/upload"
		got, err := NewBasicUploader(c, WithEndpoint(want))
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}
		if got.url != want {
			t.Errorf("want: %v, got: %v", want, got.log)
		}
	})

	t.Run("WithEmptyEndpoint", func(t *testing.T) {
		_, err := NewBasicUploader(c, WithEndpoint(""))
		if err == nil {
			t.Errorf("error was expected")
		}
	})
}

func TestBasicUploader_Upload(t *testing.T) {
	i := mockedUploadItem{}

	t.Run("ReturnsTokenOnSuccessfulUpload", func(t *testing.T) {
		want := "token"
		c := &mockedHttpClient{
			res: responseAs(http.StatusOK, "", want),
			code: http.StatusOK,
		}

		u, err := NewBasicUploader(c)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		got, err := u.Upload(context.Background(), i)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		if got != UploadToken(want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("ReturnsErrorOnFailedUpload", func(t *testing.T) {
		c := &mockedHttpClient{
			res: responseAs(http.StatusTooManyRequests, "", "token"),
			code: http.StatusTooManyRequests,
		}

		u, err := NewBasicUploader(c)
		if err != nil {
			t.Fatalf("error was not expected at this point. err: %s", err)
		}

		if _, err := u.Upload(context.Background(), i); err == nil {
			t.Errorf("error was expected.")
		}
	})
}

type mockedUploadItem struct {
	id   string
	size int64
}

func (m mockedUploadItem) Open() (io.ReadSeeker, int64, error) {
	var b bytes.Buffer
	r := strings.NewReader("some io.Reader stream to be read\n")
	size, err := b.ReadFrom(r)
	if err != nil {
		return r, 0, err
	}
	m.size = size
	return r, size, nil
}

func (m mockedUploadItem) Name() string {
	return m.String()
}

func (m mockedUploadItem) String() string {
	return m.id
}

func (m mockedUploadItem) Size() int64 {
	return m.size
}