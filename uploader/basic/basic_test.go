package basic_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/basic"
)

func TestNewBasicUploader(t *testing.T) {
	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := basic.NewBasicUploader(http.DefaultClient)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		_, err := basic.NewBasicUploader(http.DefaultClient, basic.WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := basic.NewBasicUploader(http.DefaultClient, basic.WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})
}

func TestBasicUploader_UploadFile(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		want        string
		errExpected bool
	}{
		{"Upload should be successful", "testdata/upload-success", "apiToken", false},
		{"Upload existing file with errors should be a failure", "testdata/upload-failure", "", true},
		{"Upload a non-existing file should be a failure", "non-existent", "", true},
	}
	srv := serverMock()
	defer srv.Close()

	u, err := basic.NewBasicUploader(http.DefaultClient, basic.WithEndpoint(srv.URL+"/v1/uploads"))
	if err != nil {
		t.Fatalf("error was not expected at this point, err: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := u.UploadFile(context.Background(), tc.path)
			assertExpectedError(tc.errExpected, err, t)
			if err == nil && tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func assertExpectedError(errExpected bool, err error, t *testing.T) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/v1/uploads", uploadsMock)

	return httptest.NewServer(handler)
}

func uploadsMock(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-Goog-Upload-File-Name") {
	case "upload-failure":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		_, _ = w.Write([]byte("apiToken"))
	}
}

// MockedUploadItem represents a mocked file upload item.
type MockedUploadItem struct {
	Path string
	size int64
}

// Open returns a io.ReadSeeker with a fixed string: "some test content inside a mocked file".
func (m MockedUploadItem) Open() (io.ReadSeeker, int64, error) {
	var b bytes.Buffer
	var err error

	r := strings.NewReader("some test content inside a mocked file")
	m.size, err = b.ReadFrom(r)
	if err != nil {
		return r, 0, err
	}
	return r, m.size, nil
}

// Name returns the name (path) of the item.
func (m MockedUploadItem) Name() string {
	return m.Path
}

// Size returns the length of "some test content inside a mocked file".
func (m MockedUploadItem) Size() int64 {
	return m.size
}
