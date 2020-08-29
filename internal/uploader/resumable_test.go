package uploader

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

func TestNewResumableUploader(t *testing.T) {
	c := http.DefaultClient
	s := &mockUploadSessionStore{}

	t.Run("WithoutOptions", func(t *testing.T) {
		got, err := NewResumableUploader(c, s)
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
		if got.store != s {
			t.Errorf("want: %v, got: %v", s, got.log)
		}
	})

	t.Run("WithOptionLog", func(t *testing.T) {
		want := &log.DiscardLogger{}
		got, err := NewResumableUploader(c, s, WithLogger(want))
		if err != nil {
			t.Fatalf("error was not expected here: err=%s", err)
		}
		if got.log != want {
			t.Errorf("want: %v, got: %v", want, got.log)
		}
	})

	t.Run("WithOptionEndpoint", func(t *testing.T) {
		want := "https://localhost/test/TestMe"
		got, err := NewResumableUploader(c, s, WithEndpoint(want))
		if err != nil {
			t.Errorf("NewUploader error was not expected here: err=%s", err)
		}
		if got.url != want {
			t.Errorf("want: %v, got: %v", want, got.url)
		}
	})

	t.Run("WithNilSessionStore", func(t *testing.T) {
		_, err := NewResumableUploader(c, nil)
		if err == nil {
			t.Errorf("error was expected when store in nil")
		}
	})
}

func Test_createInitialUploadRequest(t *testing.T) {
	type test struct {
		url      string
		name     string
		size     int64
		wantName string
	}

	tests := []test{
		{url: "", name: "/foo/bar/file.jpg", size: 1, wantName: "file.jpg"},
		{url: "", name: "/foo/bar/xyz.jpg", size: 1024, wantName: "xyz.jpg"},
		{url: "https://localhost/test/TestMe", name: "/foo/xyz.jpg", size: 0, wantName: "xyz.jpg"},
		{url: "https://localhost/test/TestMe", name: "file.jpg", size: 1024, wantName: "file.jpg"},
	}

	for i, tt := range tests {
		upload := &Upload{
			name: tt.name,
			size: tt.size,
		}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := upload.createInitialUploadRequest(tt.url)
			if err != nil {
				t.Fatalf("error was not expected: err=%s", err)
			}
			gotName := req.Header.Get("X-Goog-Upload-File-Name")
			if gotName != tt.wantName {
				t.Errorf("name: got=%s, want=%s", gotName, tt.wantName)
			}
			gotSize, err := strconv.ParseInt(req.Header.Get("X-Goog-Upload-Raw-Size"), 10, 64)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
			}
			if gotSize != tt.size {
				t.Errorf("maxBytes: got=%d, want=%d", gotSize, tt.size)
			}
			gotURL := req.URL.String()
			if gotURL != tt.url {
				t.Errorf("url: got=%s, want=%s", gotURL, tt.url)
			}
		})
	}
}

func Test_createQueryOffsetRequest(t *testing.T) {
	type test struct {
		url string
	}

	tests := []test{
		{url: ""},
		{url: "https://localhost/test/TestMe"},
		{url: "https://abc/"},
		{url: "https://abc/def"},
	}
	for i, tt := range tests {
		upload := &Upload{}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := upload.createQueryOffsetRequest(tt.url)
			if err != nil {
				t.Fatalf("error was not expected: err=%s", err)
			}
			got := req.URL.String()
			if got != tt.url {
				t.Errorf("url: got=%s, want=%s", got, tt.url)
			}
		})
	}
}

func Test_createResumeUploadRequest(t *testing.T) {
	type test struct {
		url    string
		size   int64
		offset int64
	}

	tests := []test{
		{url: "https://localhost/test/TestMe", size: 1, offset: 0},
		{url: "https://abc/def", size: 1024, offset: 512},
	}

	for i, tt := range tests {
		upload := &Upload{size: tt.size, sent: tt.offset}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := upload.createResumeUploadRequest(tt.url)
			if err != nil {
				t.Fatalf("error was not expected: err=%s", err)
			}
			gotOffset, err := strconv.ParseInt(req.Header.Get("X-Goog-Upload-Offset"), 10, 64)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
			}
			if gotOffset != tt.offset {
				t.Errorf("maxBytes: got=%d, want=%d", gotOffset, tt.offset)
			}
			gotURL := req.URL.String()
			if gotURL != tt.url {
				t.Errorf("url: got=%s, want=%s", gotURL, tt.url)
			}
		})
	}
}

type mockUploadSessionStore struct{}

func (m *mockUploadSessionStore) Get(f string) []byte {
	return []byte(f)
}

func (m *mockUploadSessionStore) Set(f string, u []byte) {}

func (m *mockUploadSessionStore) Delete(f string) {}
