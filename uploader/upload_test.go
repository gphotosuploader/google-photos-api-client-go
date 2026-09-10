package uploader_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"os"
	"reflect"
	"testing"
)

func TestNewUploadFromFile(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.path)
			if err != nil {
				t.Fatalf("error was not expected at this point: %s", err)
			}
			upload, err := uploader.NewUploadFromFile(f)
			if err != nil && !tc.isErrExpected {
				t.Errorf("error was not expected, err: %s", err)
			}
			if err == nil {
				data, readErr := os.ReadFile(tc.path)
				if readErr != nil {
					t.Fatal(readErr)
				}
				sum := sha256.Sum256(data)
				if upload.SHA256 != hex.EncodeToString(sum[:]) {
					t.Fatalf("SHA-256 mismatch")
				}
			}
		})
	}
}

func TestNewUploadFromFileFingerprintIncludesPath(t *testing.T) {
	dir := t.TempDir()
	for _, sub := range []string{"a", "b"} {
		if err := os.MkdirAll(dir+"/"+sub, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dir+"/"+sub+"/same.jpg", []byte("same bytes"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	open := func(path string) *uploader.Upload {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		u, err := uploader.NewUploadFromFile(f)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}
	if open(dir+"/a/same.jpg").Fingerprint == open(dir+"/b/same.jpg").Fingerprint {
		t.Fatal("different paths must not share a resumable-upload fingerprint")
	}
}

func TestNewUpload(t *testing.T) {
	input := "foo bar baz"
	r := bytes.NewBuffer([]byte(input)) // Not a reader.(io.ReadSeeker)

	wantName := "fooString"
	wantFingerprint := "fooFingerprint"

	upload := uploader.NewUpload(r, int64(reflect.TypeOf(input).Size()), wantName, wantFingerprint)

	if wantName != upload.Name {
		t.Errorf("want: %s, got: %s", wantName, upload.Name)
	}

	if wantFingerprint != upload.Fingerprint {
		t.Errorf("want: %s, got: %s", wantFingerprint, upload.Fingerprint)
	}
}
