package uploader

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
)

func Test_createInitialResumableUploadRequest(t *testing.T) {
	type test struct {
		url  string
		size int64
	}

	tests := []test{
		{url: "", size: 1},
		{url: "", size: 1024},
		{url: "https://localhost/test/TestMe", size: 0},
		{url: "https://localhost/test/TestMe", size: 1024},
	}

	for i, tt := range tests {
		upload := &Upload{size: tt.size}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := createInitialResumableUploadRequest(tt.url, upload)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
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

func Test_createRawUploadRequest(t *testing.T) {
	type test struct {
		url  string
		name string
	}

	tests := []test{
		{url: "", name: "testTest"},
		{url: "https://localhost/test/TestMe", name: "testTest"},
		{url: "https://localhost/test/TestMe", name: "testTest"},
	}

	l := log.New(ioutil.Discard, "", 0)
	for i, tt := range tests {
		upload := &Upload{name: tt.name}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := createRawUploadRequest(tt.url, upload, l)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
			}
			gotName := req.Header.Get("X-Goog-Upload-File-Name")
			if gotName != tt.name {
				t.Errorf("name: got=%s, want=%s", gotName, tt.name)
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
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := createQueryOffsetRequest(tt.url)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
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

	l := log.New(ioutil.Discard, "", 0)
	for i, tt := range tests {
		upload := &Upload{size: tt.size, sent: tt.offset}
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			req, err := createResumeUploadRequest(tt.url, upload, l)
			if err != nil {
				t.Errorf("error was not expected: err=%s", err)
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
