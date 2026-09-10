package uploader

import (
	"io"
	"time"
)

type rateLimitedReader struct {
	r              io.Reader
	bytesPerSecond int64
	started        time.Time
	read           int64
}

func newRateLimitedReader(r io.Reader, bytesPerSecond int64) io.Reader {
	if bytesPerSecond <= 0 {
		return r
	}
	return &rateLimitedReader{r: r, bytesPerSecond: bytesPerSecond, started: time.Now()}
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.read += int64(n)
	target := time.Duration(float64(r.read) / float64(r.bytesPerSecond) * float64(time.Second))
	if wait := target - time.Since(r.started); wait > 0 {
		time.Sleep(wait)
	}
	return n, err
}
