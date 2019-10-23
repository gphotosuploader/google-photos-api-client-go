package log

func DefaultLogger() Logger {
	return &DiscardLogger{}
}
