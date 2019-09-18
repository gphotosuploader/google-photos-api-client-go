package gphotos

import (
	"fmt"
	"io"
	"log"
	"os"
)

// ReadProgressReporter represents read progress.
type ReadProgressReporter struct {
	r      io.Reader   // where to read data from.
	logger *log.Logger // where to log progress status.

	filename string // name of the file being uploaded
	size     int64  // size of the file
	sent     int64  // bytes already sent
	atEOF    bool   // file has reach EOF
}

func DefaultReadProgressReporter(r io.Reader, filename string, size, sent int64) ReadProgressReporter {
	return ReadProgressReporter{
		r:      r,
		logger: log.New(os.Stdout, LogPrefix, log.LstdFlags),

		filename: filename,
		size:     size,
		sent:     sent,
		atEOF:    false,
	}
}

func (pr *ReadProgressReporter) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.sent += int64(n)
	if err == io.EOF {
		pr.atEOF = true
	}
	pr.report()
	return n, err
}

func (pr *ReadProgressReporter) report() {
	if pr.atEOF {
		pr.logger.Printf("Upload completed: file=%s", pr.filename)
		return
	}

	pr.logger.Print(pr.progressLine())
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
	return fmt.Sprintf("[%s] Sent %d of %d bytes: file=%s", pr.completedPercentString(), pr.sent, pr.size, pr.filename)
}
