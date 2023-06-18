package uploader_test

import (
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"os"
	"testing"
)

var testCases = []struct {
	name          string
	path          string
	wantName      string
	wantSize      int64
	isErrExpected bool
}{
	{name: "sample JPEG 100kB", path: "testdata/file_example_JPG_100kB.jpg", wantName: "file_example_JPG_100kB.jpg", wantSize: 102117, isErrExpected: false},
	{name: "sample PNG 500kB", path: "testdata/file_example_PNG_500kB.png", wantName: "file_example_PNG_500kB.png", wantSize: 512596, isErrExpected: false},
	{name: "sample WEBP 50kB", path: "testdata/file_example_WEBP_50kB.webp", wantName: "file_example_WEBP_50kB.webp", wantSize: 50408, isErrExpected: false},
}

func TestNewUploadFromFile(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.path)
			if err != nil {
				t.Fatalf("error was not expected at this point: %s", err)
			}
			_, err = uploader.NewUploadFromFile(f)
			if err != nil && !tc.isErrExpected {
				t.Errorf("error was not expected, err: %s", err)
			}
		})
	}
}
