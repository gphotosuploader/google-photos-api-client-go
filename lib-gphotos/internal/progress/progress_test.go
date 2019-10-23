package progress_test

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos/internal/log"
)

func TestReadProgressReporter_completedPercentString(t *testing.T) {
	var testData = []struct {
		sent, size int64
		want       string
	}{
		{sent: 0, size: 100, want: "0%"},
		{sent: 5, size: 100, want: "5%"},
		{sent: 15, size: 100, want: "15%"},
		{sent: 100, size: 100, want: "100%"},
		{sent: 0, size: 0, want: "N/A"},
		{sent: 100, size: 0, want: "N/A"},
		{sent: -1, size: 100, want: "N/A"},
		{sent: math.MaxInt64 / 2, size: math.MaxInt64, want: "50%"},
	}

	l := &log.DiscardLogger{}
	for tn, tt := range testData {
		r := Reporter{
			reader:       nil,
			logger:       l,
			description:  "testTest",
			maxBytes:     tt.size,
			currentBytes: tt.sent,
		}
		got := r.completedPercentString()
		if got != tt.want {
			t.Errorf("test number %d failed: got=%s, want=%s", tn+1, got, tt.want)
		}
	}
}

func TestReadProgressReporter_report(t *testing.T) {
	var testData = []struct {
		finished bool
		want     string
	}{
		{
			finished: false,
			want:     "[10%] Sent 10 of 100 bytes: file=testTest",
		},
		{
			finished: true,
			want:     "Upload completed: file=testTest",
		},
	}

	for tn, tt := range testData {
		var buff = bytes.Buffer{}
		logger :=

		r := Reporter{
			reader: nil,
			logger: logger,

			description:  "testTest",
			maxBytes:     100,
			currentBytes: 10,
			finished:     tt.finished,
		}

		r.render()
		got := strings.TrimSuffix(buff.String(), "\n")
		if got != tt.want {
			t.Errorf("test number %d failed: got=%s, want=%s", tn+1, got, tt.want)
		}
		buff.Truncate(0)
	}
}

func TestReadProgressReporter_Read(t *testing.T) {
	want := "abcde"
	r := Reporter{
		reader: strings.NewReader(want),
		logger: log.New(&bytes.Buffer{}, "", 0),

		description:  "testTest",
		maxBytes:     int64(len(want)),
		currentBytes: 0,
		finished:     false,
	}

	buf := make([]byte, len(want))
	_, err := r.Read(buf)
	if err != nil {
		t.Errorf("error was not expected: err=%s", err)
	}

	got := string(buf)
	if got != want {
		t.Errorf("test failed: got=%s, want=%s", got, want)
	}
}
