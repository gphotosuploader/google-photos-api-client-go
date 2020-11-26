package resumable_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/resumable"
)

func TestNewResumableUploader(t *testing.T) {
	s := &MockedSessionStorer{}

	t.Run("WithoutOptions", func(t *testing.T) {
		_, err := resumable.NewResumableUploader(http.DefaultClient, s)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &MockedLogger{}
		_, err := resumable.NewResumableUploader(http.DefaultClient, s, resumable.WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		_, err := resumable.NewResumableUploader(http.DefaultClient, s, resumable.WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
	})

	t.Run("WithNilSessionStore", func(t *testing.T) {
		_, err := resumable.NewResumableUploader(http.DefaultClient, nil)
		if err == nil {
			t.Errorf("error was expected when store in nil")
		}
	})
}

func TestResumableUploader_UploadFile(t *testing.T) {
	testCases := []struct {
		name            string
		path            string
		alreadyUploaded bool
		want            string
		errExpected     bool
	}{
		{"Should be successful when file is uploaded", "testdata/upload-success", false, "apiToken", false},
		{"Should fail when file is not uploaded", "testdata/upload-should-fail", false, "", true},
		//{"Upload a non-existing file should be a failure", "non-existent", false, "", true},
	}
	srv := NewMockedGooglePhotosServer()
	defer srv.Close()

	s := &MockedSessionStorer{
		GetFn: func(f string) []byte {
			return []byte(srv.URL("/upload-session"))
		},
		SetFn:    func(f string, u []byte) {},
		DeleteFn: func(f string) {},
	}
	l := &MockedLogger{
		LogFn: func(args ...interface{}) {
			fmt.Println(args...)
		},
	}

	u, err := resumable.NewResumableUploader(http.DefaultClient, s, resumable.WithEndpoint(srv.URL("/v1/uploads")), resumable.WithLogger(l))
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

// MockedGooglePhotosServer mock the Google Photos Service for uploads.
type MockedGooglePhotosServer struct {
	server  *httptest.Server
	baseURL string
}

func NewMockedGooglePhotosServer() *MockedGooglePhotosServer {
	ms := &MockedGooglePhotosServer{}
	ms.server = httptest.NewServer(ms.routes())
	ms.baseURL = ms.server.URL
	return ms
}

func (ms MockedGooglePhotosServer) routes() *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("/v1/uploads", ms.handleUploads)
	handler.HandleFunc("/upload-session", ms.handleUploadSession)
	return handler
}

func (ms MockedGooglePhotosServer) URL(endpoint string) string {
	return ms.baseURL + endpoint
}

func (ms MockedGooglePhotosServer) Close() {
	ms.server.Close()
}

func (ms MockedGooglePhotosServer) handleUploadSession(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-Goog-Upload-Command") {
	case "upload, finalize":
		// get: X-Goog-Upload-Offset
		_, _ = w.Write([]byte("apiToken"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (ms MockedGooglePhotosServer) handleUploads(w http.ResponseWriter, r *http.Request) {
	if "upload-should-fail" == r.Header.Get("X-Goog-Upload-File-Name") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Header.Get("X-Goog-Upload-Command") {
	case "start":
		// get: X-Goog-Upload-File-Name
		// get: X-Goog-Upload-Raw-Size
		w.Header().Add("X-Goog-Upload-URL", ms.URL("/upload-session"))
	case "query":
		w.Header().Add("X-Goog-Upload-Status", "active") // other values could be: final and cancelled
		w.Header().Add("X-Goog-Upload-Size-Received", "0")
	case "upload, finalize":
		// get: X-Goog-Upload-Offset
		_, _ = w.Write([]byte("apiToken"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// MockedSessionStorer mocks a service to store resumable upload data.
type MockedSessionStorer struct {
	GetFn    func(f string) []byte
	SetFn    func(f string, u []byte)
	DeleteFn func(f string)
}

func (s MockedSessionStorer) Get(f string) []byte {
	return s.GetFn(f)
}

func (s MockedSessionStorer) Set(f string, u []byte) {
	s.SetFn(f, u)
}

func (s MockedSessionStorer) Delete(f string) {
	s.DeleteFn(f)
}

// MockedLogger mocks a logger.
type MockedLogger struct {
	LogFn func(args ...interface{})
}

func (d *MockedLogger) Debug(args ...interface{}) {
	d.LogFn(args...)
}

func (d *MockedLogger) Debugf(format string, args ...interface{}) {
	d.LogFn(fmt.Sprintf(format, args...))
}

func (d *MockedLogger) Info(args ...interface{}) {
	d.LogFn(args...)
}

func (d *MockedLogger) Infof(format string, args ...interface{}) {
	d.LogFn(fmt.Sprintf(format, args...))
}

func (d *MockedLogger) Warn(args ...interface{}) {
	d.LogFn(args...)
}

func (d *MockedLogger) Warnf(format string, args ...interface{}) {
	d.LogFn(fmt.Sprintf(format, args...))
}

func (d *MockedLogger) Error(args ...interface{}) {
	d.LogFn(args...)
}

func (d *MockedLogger) Errorf(format string, args ...interface{}) {
	d.LogFn(fmt.Sprintf(format, args...))
}
