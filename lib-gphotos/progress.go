package gphotos

import (
	"fmt"
	"io"
	"os"
)

// ReadProgressReporter represents read progress.
type ReadProgressReporter struct {
	r   io.Reader // where to read data from.
	out io.Writer // where to write progress status.

	size  int64 // size of the file
	sent  int64 // bytes already sent
	atEOF bool  // file has reach EOF
}

func DefaultReadProgressReporter(r io.Reader, size, sent int64) ReadProgressReporter {
	return ReadProgressReporter{
		r:     r,
		out:   os.Stdout,
		size:  size,
		sent:  sent,
		atEOF: false,
	}
}

func (pr *ReadProgressReporter) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.sent += int64(n)
	if err == io.EOF {
		pr.atEOF = true
	}
	err = pr.report()
	return n, err
}

func (pr *ReadProgressReporter) report() error {
	var statusCompleted string

	if pr.atEOF {
		statusCompleted = " Upload completed."
	}

	_, err := fmt.Fprintf(pr.out, "%s", pr.progressLine()+statusCompleted)
	return err
}

// completedPercent return the percent completed.
func (pr *ReadProgressReporter) completedPercent() int {
	if pr.size <= 0 || pr.sent < 0 {
		return -1
	}
	completed := float64(pr.sent) / float64(pr.size)
	return int(completed * 100)
}

// completedPercentString returns the formatted string representation of the completed percent
func (pr *ReadProgressReporter) completedPercentString() string {
	cp := pr.completedPercent()
	if cp < 0 {
		return "N/A"
	}
	return fmt.Sprintf("%d%%", cp)
}

func (pr *ReadProgressReporter) progressLine() string {
	return fmt.Sprintf("[%s] Sent %d of %d bytes.", pr.completedPercentString(), pr.sent, pr.size)
}
