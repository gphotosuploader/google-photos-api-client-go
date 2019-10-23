package progress

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Reporter represents io.Reader with progress render.
type Reporter struct {
	reader io.Reader // where to read data from.
	writer io.Writer // where to write progress status.

	description string // name or description of the reader
	max         int64  // maximum size of bytes that can be read from reader
	current     int64  // bytes already read from reader
	finished    bool   // reader has reach EOF
}

// NewOptions constructs a new instance of progress Reporter, with any options you specify
func NewOptions(r io.Reader, max int, options ...Option) *Reporter {
	return NewOptions64(r, int64(max), options...)
}

// NewOptions64 constructs a new instance of progress Reporter, with any options you specify
func NewOptions64(r io.Reader, max int64, options ...Option) *Reporter {
	rp := Reporter{
		reader: r,
		writer: os.Stdout,
		max:    max,
	}

	for _, o := range options {
		o(&rp)
	}

	return &rp
}

// New returns a new progress Reporter with the specified maximum
func New(r io.Reader, max int) *Reporter {
	return NewOptions(r, max)
}

// Option is the type all options need to adhere to
type Option func(rp *Reporter)

// OptionSetWriter sets the output writer (defaults to os.StdOut)
func WithWriter(w io.Writer) Option {
	return func(rp *Reporter) {
		rp.writer = w
	}
}

// OptionSetDescription sets the description of the progress render to render in front of it
func WithDescription(description string) Option {
	return func(rp *Reporter) {
		rp.description = description
	}
}

// Read will read the data and add the number of bytes to the progress Reporter
func (pr *Reporter) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	_ = pr.Add(n)
	return n, err
}

// add will add the specified amount to the progress Reporter
func (pr *Reporter) add(num int) error {
	return pr.Add64(int64(num))
}

// add64 will add the specified amount to the progress Reporter
func (pr *Reporter) add64(num int64) error {
	if pr.max <= 0 {
		return errors.New("max must be greater than 0")
	}
	pr.current += num
	if pr.current > pr.max {
		return errors.New("current number exceeds max")
	}
	return pr.render()
}

// render renders the progress reporter
func (pr *Reporter) render() error {
	// check if the progress bar is finished
	if !pr.finished && pr.current >= pr.max {
		pr.finished = true
	}
	if pr.finished {
		return nil
	}

	// then, re-render the current progress reporter
	_, err := io.WriteString(pr.writer, pr.progressLine())
	return err
}

// completedPercent return the percent completed.
func (pr *Reporter) completedPercent() int {
	if pr.max <= 0 || pr.current < 0 {
		return -1
	}
	percent := float64(pr.current) / float64(pr.max)
	return int(percent * 100)
}

// completedPercentString returns the formatted string representation of the completed percent
func (pr *Reporter) completedPercentString() string {
	cp := pr.completedPercent()
	if cp < 0 {
		return "N/A"
	}
	return fmt.Sprintf("%d%%", cp)
}

func (pr *Reporter) progressLine() string {
	pl := fmt.Sprintf("[%s] Sent %d of %d bytes", pr.completedPercentString(), pr.current, pr.max)
	if pr.description != "" {
		pl += ": file=" + pr.description
	}
	return pl
}
