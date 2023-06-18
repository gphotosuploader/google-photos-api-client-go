package uploader_test

import (
	"context"
	"fmt"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/mocks"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"net/http"
	"testing"
)

func TestNewResumableUploader(t *testing.T) {
	u, err := uploader.NewResumableUploader(http.DefaultClient)
	if err != nil {
		t.Fatalf("error was not expected at this point: %s", err)
	}
	want := "https://photoslibrary.googleapis.com/v1/uploads"

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
		{"Should be successful when file is resuming createOrResumeUpload ", "testdata/upload-resume-success", true, "apiToken", false},
		{"Should fail when file is not uploaded", "testdata/upload-should-fail", false, "", true},
		{"Should fail if file doesn't exist", "non-existent", false, "", true},
	}
	srv := mocks.NewMockedGooglePhotosService()
	defer srv.Close()

	store := NewMockStore()

	logger := &MockedLogger{
		LogFn: func(args ...interface{}) {
			//fmt.Fprintln(os.Stderr, args...)
		},
	}

	u, err := uploader.NewResumableUploader(http.DefaultClient)
	u.BaseURL = srv.URL() + "/v1/uploads"
	u.Store = store
	u.Logger = logger

	if err != nil {
		t.Fatalf("error was not expected at this point, err: %s", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := u.UploadFile(context.Background(), tc.path)
			if tc.errExpected && err == nil {
				t.Fatalf("error was expected, but not produced")
			}
			if !tc.errExpected && err != nil {
				t.Fatalf("error was not expected, err: %s", err)
			}
			if err == nil && uploader.UploadToken(tc.want) != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

type MockStore struct {
	m map[string]string
}

func NewMockStore() uploader.Store {
	return &MockStore{
		make(map[string]string),
	}
}

func (s *MockStore) Get(fingerprint string) (string, bool) {
	url, ok := s.m[fingerprint]
	return url, ok
}

func (s *MockStore) Set(fingerprint, url string) {
	s.m[fingerprint] = url
}

func (s *MockStore) Delete(fingerprint string) {
	delete(s.m, fingerprint)
}

func (s *MockStore) Close() {
	for k := range s.m {
		delete(s.m, k)
	}
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
