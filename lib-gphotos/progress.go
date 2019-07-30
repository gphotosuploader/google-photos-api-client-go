package gphotos

import (
	"fmt"
	"io"
)

type ReadProgressReporter struct {
	r        io.Reader
	max      int
	sent     int
	atEOF    bool
	fileSize int
}

func (pr *ReadProgressReporter) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.sent += n
	if err == io.EOF {
		pr.atEOF = true
	}
	pr.report()
	return n, err
}

func (pr *ReadProgressReporter) report() {
	percent := (pr.fileSize - pr.max + pr.sent) * 100 / pr.fileSize
	fmt.Printf("\r[%d%%] Sent %d of %d bytes (total file size: %d)", percent, pr.sent, pr.max, pr.fileSize)
	if pr.atEOF {
		fmt.Println("\nUpload done")
	}
}
