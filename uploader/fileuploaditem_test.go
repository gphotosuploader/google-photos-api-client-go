package uploader

import (
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
	{name: "non existent file", path: "testdata/non-existent-file.txt", wantName: "", wantSize: 0, isErrExpected: true},
}

func TestFileUploadItem_Open(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := FileUploadItem(tc.path)
			if _, _, err := f.Open(); err != nil && !tc.isErrExpected {
				t.Errorf("error was not expected, err: %s", err)
			}
		})
	}
}

func TestFileUploadItem_Name(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := FileUploadItem(tc.path)
			if got := f.Name(); tc.wantName != got {
				t.Errorf("want: %s, got: %s", tc.wantName, got)
			}
		})
	}
}

func TestFileUploadItem_Size(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := FileUploadItem(tc.path)
			if got := f.Size(); tc.wantSize != got {
				t.Errorf("want: %d, got: %d", tc.wantSize, got)
			}
		})
	}
}
