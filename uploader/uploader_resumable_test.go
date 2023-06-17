package gphotos_test

import (
	"context"
	"fmt"
	"github.com/gphotosuploader/google-photos-api-client-go/v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResumableUploader(t *testing.T) {
	u, err := gphotos.NewResumableUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	want := gphotos.DefaultEndpoint
	if want != u.BaseURL {
		t.Errorf("want: %s, got: %s", want, u.BaseURL)
	}
}

func TestResumableUploader_UploadFile(t *testing.T) {
	testCases := []struct {
		name           string
		path           string
		alreadyStarted bool
		want           string
		errExpected    bool
	}{
		{"Should be successful when file is uploaded", "testdata/upload-success", false, "apiToken", false},
		{"Should be successful when file is resuming upload ", "testdata/upload-resume-success", true, "apiToken", false},
		{"Should fail when file is not uploaded", "testdata/upload-should-fail", false, "", true},
		{"Should fail if file doesn't exist", "non-existent", false, "", true},
	}
	srv := NewMockedGooglePhotosServer()
	defer srv.Close()

	u, err := gphotos.NewResumableUploader(http.DefaultClient)
	u.BaseURL = srv.URL("/uploads")

	if err != nil {
		t.Fatalf("error was not expected at this point, err: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//sessionStorerData = []byte{}
			if tc.alreadyStarted {
				//sessionStorerData = []byte(srv.URL("/upload-session/started"))
			}
			got, err := u.UploadFile(context.Background(), tc.path)
			assertExpectedError(tc.errExpected, err, t)
			if err == nil && gphotos.UploadToken(tc.want) != got {
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
	mux := http.NewServeMux()
	ms.server = httptest.NewServer(mux)
	ms.baseURL = ms.server.URL
	mux.HandleFunc("/uploads", ms.handleUploads)
	mux.HandleFunc("/upload-session/started", ms.handleExistingUploadSession)
	mux.HandleFunc("/upload-session/upload-success", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("apiToken"))
	})
	return ms
}

func (ms MockedGooglePhotosServer) URL(endpoint string) string {
	return ms.baseURL + endpoint
}

func (ms MockedGooglePhotosServer) Close() {
	ms.server.Close()
}

func (ms MockedGooglePhotosServer) handleExistingUploadSession(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-Goog-Upload-Command") {
	case "query":
		w.Header().Add("X-Goog-Upload-Status", "active")
		w.Header().Add("X-Goog-Upload-Size-Received", "1000")
	case "upload, finalize":
		_, _ = w.Write([]byte("apiToken"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (ms MockedGooglePhotosServer) handleUploads(w http.ResponseWriter, r *http.Request) {
	if "start" != r.Header.Get("X-Goog-Upload-Command") {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Header.Get("X-Goog-Upload-File-Name") {
	case "upload-should-fail":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.Header().Add("X-Goog-Upload-URL", ms.URL("/upload-session/upload-success"))
	}
}

// MockedSessionStorer mocks a service to Store resumable upload data.
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
