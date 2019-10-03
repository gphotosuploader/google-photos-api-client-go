package uploader

import (
	"fmt"
	"io"
)

func (u *Upload) Read(p []byte) (int, error) {
	n, err := u.r.Read(p)
	u.sent += int64(n)
	if err == io.EOF {
		u.atEOF = true
	}
	u.report()
	return n, err
}

func (u *Upload) report() {
	if u.atEOF {
		u.Uploader.log.Printf("Upload completed: file=%s", u.name)
		return
	}

	u.Uploader.log.Print(u.progressLine())
}

// completedPercent return the percent completed.
func (u *Upload) completedPercent() int {
	if u.size <= 0 || u.sent < 0 {
		return -1
	}
	completed := float64(u.sent) / float64(u.size)
	return int(completed * 100)
}

// completedPercentString returns the formatted string representation of the completed percent
func (u *Upload) completedPercentString() string {
	cp := u.completedPercent()
	if cp < 0 {
		return "N/A"
	}
	return fmt.Sprintf("%d%%", cp)
}

func (u *Upload) progressLine() string {
	return fmt.Sprintf("[%s] Sent %d of %d bytes: file=%s", u.completedPercentString(), u.sent, u.size, u.name)
}
