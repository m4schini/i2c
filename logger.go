package i2c

var _logger = NopLogger

func NopLogger(format string, args ...any) {}

func UseLogger(f func(format string, args ...any)) {
	if f == nil {
		_logger("Aborted: New logger is nil. If you want to disable logs pass NopLogger")
		return
	}
	f("Using new logger")
	_logger = f
}

func logf(format string, args ...any) {
	_logger(format, args...)
}
