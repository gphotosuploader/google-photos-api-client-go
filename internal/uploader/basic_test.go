package uploader

import (
	"fmt"
	"testing"
)

func Test_createRawUploadRequest(t *testing.T) {
	type test struct {
		url      string
		name     string
		wantName string
	}

	tests := []test{
		{url: "", name: "/foo/xyz.jpg", wantName: "xyz.jpg"},
		{url: "https://localhost/test/TestMe", name: "/foo/bar/file.jpg", wantName: "file.jpg"},
		{url: "https://localhost/test/TestMe", name: "/foo/xyz.jpg", wantName: "xyz.jpg"},
	}

	for i, tt := range tests {
		upload := &Upload{name: tt.name}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := upload.createRawUploadRequest(tt.url)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
			}
			gotName := req.Header.Get("X-Goog-Upload-File-Name")
			if gotName != tt.wantName {
				t.Errorf("name: got=%s, want=%s", gotName, tt.wantName)
			}
			gotURL := req.URL.String()
			if gotURL != tt.url {
				t.Errorf("url: got=%s, want=%s", gotURL, tt.url)
			}
		})
	}
}
