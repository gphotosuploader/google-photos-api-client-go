package uploader

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestRateLimiterPreservesEveryByte(t *testing.T) {
	source := bytes.Repeat([]byte("original-quality-bytes"), 2048)
	started := time.Now()
	got, err := io.ReadAll(newRateLimitedReader(bytes.NewReader(source), int64(len(source))*10))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, source) {
		t.Fatal("rate limiter changed uploaded bytes")
	}
	if time.Since(started) < 80*time.Millisecond {
		t.Fatal("rate limiter did not pace the stream")
	}
}
