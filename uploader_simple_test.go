package gphotos_test

import (
	"context"
	"github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestNewBasicUploader(t *testing.T) {
	got, err := gphotos.NewSimpleUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	want := uploader.DefaultEndpoint

	if want != got.BaseURL {
		t.Errorf("want: %s, got: %s", want, got)
	}
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
	srv := NewMockedGooglePhotosServer()
	defer srv.Close()

	u, err := gphotos.NewSimpleUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	u.BaseURL = srv.URL("/uploads")

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
	mux := http.NewServeMux()
	ms.server = httptest.NewServer(mux)
	ms.baseURL = ms.server.URL
	mux.HandleFunc("/uploads", ms.handleUploads)
	return ms
}

func (ms MockedGooglePhotosServer) Close() {
	ms.server.Close()
}

func (ms MockedGooglePhotosServer) URL(endpoint string) string {
	return ms.baseURL + endpoint
}

func (ms MockedGooglePhotosServer) handleUploads(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("X-Goog-Upload-File-Name") {
	case "upload-failure":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		var bodyContent []byte
		bodyLength, err := r.Body.Read(bodyContent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		expectedLength, _ := strconv.Atoi(r.Header.Get("Content-Length"))
		if expectedLength != bodyLength {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("apiToken"))
	}
}
