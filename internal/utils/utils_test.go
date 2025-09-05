package utils

import (
	"bytes"
	"io"
	"log"
	"testing"
)

type closerStub struct {
	closeErr error
}

func (c *closerStub) Close() error {
	return c.closeErr
}

func TestClosesWithoutLoggingWhenNoError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	c := &closerStub{closeErr: nil}
	CloseOrLog(c, "resourceA")

	if buf.Len() != 0 {
		t.Errorf("log output was not expected, got: %s", buf.String())
	}
}

func TestLogsErrorWhenCloseReturnsError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	c := &closerStub{closeErr: io.ErrUnexpectedEOF}
	CloseOrLog(c, "resourceB")

	if !bytes.Contains(buf.Bytes(), []byte("Error while closing resource \"resourceB\"")) {
		t.Errorf("expected log output to contain error message, got: %s", buf.String())
	}
}

func TestHandlesEmptyResourceNameGracefully(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	c := &closerStub{closeErr: io.ErrUnexpectedEOF}
	CloseOrLog(c, "")

	if !bytes.Contains(buf.Bytes(), []byte("Error while closing resource \"\"")) {
		t.Errorf("expected log output for empty resource name, got: %s", buf.String())
	}
}
