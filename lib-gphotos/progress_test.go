package gphotos

import (
	"bytes"
	"math"
	"testing"
)

func TestReadProgressReporter_completedPercentString(t *testing.T) {
	var testData = []struct {
		sent, size int64
		want            string
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
		r := DefaultReadProgressReporter(nil, tt.size, tt.sent)
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
			want:  "[10%] Sent 10 of 100 bytes.",
		},
		{
			atEOF: true,
			want:  "[10%] Sent 10 of 100 bytes. Upload completed.",
		},
	}

	for tn, tt := range testData {
		var buffer = &bytes.Buffer{}
		r := ReadProgressReporter{
			r:     nil,
			out:   buffer,
			sent:  10,
			atEOF: tt.atEOF,
			size:  100,
		}

		err := r.report()
		if err != nil {
			t.Errorf("error was not expected at this stage: err=%s", err)
		}

		got := buffer.String()
		if got != tt.want {
			t.Errorf("test number %d failed: got=%s, want=%s", tn+1, got, tt.want)
		}
		buffer.Truncate(0)
	}

}
