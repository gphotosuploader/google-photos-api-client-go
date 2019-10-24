package log

// DiscardLogger just discards every log statement
type DiscardLogger struct{}

func NewDiscardLogger() *DiscardLogger {
	return &DiscardLogger{}
}

// Debug implements logger interface
func (d *DiscardLogger) Debug(args ...interface{}) {}

// Debugf implements logger interface
func (d *DiscardLogger) Debugf(format string, args ...interface{}) {}

// Info implements logger interface
func (d *DiscardLogger) Info(args ...interface{}) {}

// Infof implements logger interface
func (d *DiscardLogger) Infof(format string, args ...interface{}) {}

// Warn implements logger interface
func (d *DiscardLogger) Warn(args ...interface{}) {}

// Warnf implements logger interface
func (d *DiscardLogger) Warnf(format string, args ...interface{}) {}

// Error implements logger interface
func (d *DiscardLogger) Error(args ...interface{}) {}

// Errorf implements logger interface
func (d *DiscardLogger) Errorf(format string, args ...interface{}) {}
