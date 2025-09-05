package log

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestLogsNothingForDebug(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Debug("This is a debug message")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForDebugf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Debugf("This is a %s message", "debug")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Info("This is an info message")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForInfof(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Infof("This is an %s message", "info")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForWarn(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Warn("This is a warn message")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForWarnf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Warnf("This is a %s message", "warn")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Error("This is an error message")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}

func TestLogsNothingForErrorf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	logger := &DiscardLogger{}
	logger.Errorf("This is an %s message", "error")

	if buf.Len() != 0 {
		t.Errorf("Expected no log output, got: %s", buf.String())
	}
}
