package gphotos

import (
	"fmt"
	"io"
)

type ReadProgressReporter struct {
	r     io.Reader
	max   int
	sent  int
	atEOF bool
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
	fmt.Printf("Sent %d of %d bytes\n", pr.sent, pr.max)
	if pr.atEOF {
		fmt.Println("DONE")
	}
}
