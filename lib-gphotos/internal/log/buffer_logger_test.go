package log

import (
	"bytes"
	"fmt"
)

// Level type
type logFunctionType uint32

const (
	panicFn logFunctionType = iota
	fatalFn
	errorFn
	warnFn
	infoFn
	debugFn
	failFn
	doneFn
)

type fnTypeInformation string

var fnTypeInformationMap = map[logFunctionType]fnTypeInformation{
	debugFn: "[debug]  ",
	infoFn:  "[info]   ",
	warnFn:  "[warn]   ",
	errorFn: "[error]  ",
	fatalFn: "[fatal]  ",
	panicFn: "[panic]  ",
	doneFn:  "[done] √ ",
	failFn:  "[fail] X ",
}

// BufferLogger just discards every log statement
type BufferLogger struct {
	buf bytes.Buffer
}

func (d *BufferLogger) writeMessage(fnType logFunctionType, message string) {
	fnInformation := fnTypeInformationMap[fnType]
	d.buf.Write([]byte(fnInformation))
	d.buf.Write([]byte(message))
}

// Debug implements logger interface
func (d *BufferLogger) Debug(args ...interface{}) {
	d.writeMessage(debugFn, fmt.Sprintln(args...))
}

// Debugf implements logger interface
func (d *BufferLogger) Debugf(format string, args ...interface{}) {
	d.writeMessage(debugFn, fmt.Sprintf(format, args...))
}

// Info implements logger interface
func (d *BufferLogger) Info(args ...interface{}) {
	d.writeMessage(infoFn, fmt.Sprintln(args...))
}

// Infof implements logger interface
func (d *BufferLogger) Infof(format string, args ...interface{}) {
	d.writeMessage(infoFn, fmt.Sprintf(format, args...))
}

// Warn implements logger interface
func (d *BufferLogger) Warn(args ...interface{}) {
	d.writeMessage(warnFn, fmt.Sprintln(args...))
}

// Warnf implements logger interface
func (d *BufferLogger) Warnf(format string, args ...interface{}) {
	d.writeMessage(warnFn, fmt.Sprintf(format, args...))
}

// Error implements logger interface
func (d *BufferLogger) Error(args ...interface{}) {
	d.writeMessage(errorFn, fmt.Sprintln(args...))
}

// Errorf implements logger interface
func (d *BufferLogger) Errorf(format string, args ...interface{}) {
	d.writeMessage(errorFn, fmt.Sprintf(format, args...))
}

// Fatal implements logger interface
func (d *BufferLogger) Fatal(args ...interface{}) {
	d.writeMessage(fatalFn, fmt.Sprintln(args...))
}

// Fatalf implements logger interface
func (d *BufferLogger) Fatalf(format string, args ...interface{}) {
	d.writeMessage(fatalFn, fmt.Sprintf(format, args...))
}

// Panic implements logger interface
func (d *BufferLogger) Panic(args ...interface{}) {
	d.writeMessage(panicFn, fmt.Sprintln(args...))
}

// Panicf implements logger interface
func (d *BufferLogger) Panicf(format string, args ...interface{}) {
	d.writeMessage(panicFn, fmt.Sprintf(format, args...))
}

// Done implements logger interface
func (d *BufferLogger) Done(args ...interface{}) {
	d.writeMessage(doneFn, fmt.Sprintln(args...))
}

// Donef implements logger interface
func (d *BufferLogger) Donef(format string, args ...interface{}) {
	d.writeMessage(doneFn, fmt.Sprintf(format, args...))
}

// Fail implements logger interface
func (d *BufferLogger) Fail(args ...interface{}) {
	d.writeMessage(failFn, fmt.Sprintln(args...))
}

// Failf implements logger interface
func (d *BufferLogger) Failf(format string, args ...interface{}) {
	d.writeMessage(failFn, fmt.Sprintf(format, args...))
}
