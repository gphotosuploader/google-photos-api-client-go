package gphotos

import (
	"bytes"
	"log"
	"math"
	"strings"
	"testing"
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

	for tn, tt := range testData {
		r := DefaultReadProgressReporter(nil, "testTest", tt.size, tt.sent)
		got := r.completedPercentString()
		if got != tt.want {
			t.Errorf("test number %d failed: got=%s, want=%s", tn+1, got, tt.want)
		}
	}
}

func TestReadProgressReporter_report(t *testing.T) {
	var testData = []struct {
		atEOF bool
		want  string
	}{
		{
			atEOF: false,
			want:  "[10%] Sent 10 of 100 bytes: file=testTest",
		},
		{
			atEOF: true,
			want:  "Upload completed: file=testTest",
		},
	}

	for tn, tt := range testData {
		var buff = bytes.Buffer{}
		logger := log.New(&buff, "", 0)

		r := ReadProgressReporter{
			r:      nil,
			logger: logger,

			filename: "testTest",
			size:     100,
			sent:     10,
			atEOF:    tt.atEOF,
		}

		r.report()
		got := strings.TrimSuffix(buff.String(), "\n")
		if got != tt.want {
			t.Errorf("test number %d failed: got=%s, want=%s", tn+1, got, tt.want)
		}
		buff.Truncate(0)
	}
}

func TestReadProgressReporter_Read(t *testing.T) {
	want := "abcde"
	r := ReadProgressReporter{
		r:      strings.NewReader(want),
		logger: log.New(&bytes.Buffer{}, "", 0),

		filename: "testTest",
		size:     int64(len(want)),
		sent:     0,
		atEOF:    false,
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
