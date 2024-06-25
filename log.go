package olog

import (
	"sync/atomic"
	"unsafe"
)

// def stores a pointer to the default logger instance.
var def unsafe.Pointer

// init sets the default logger to a new logger instance.
func init() {
	def = unsafe.Pointer(newLogger())
}

// getDefLogger returns a pointer to the default logger instance.
func getDefLogger() *logger {
	return (*logger)(atomic.LoadPointer(&def))
}

// setDefLogger sets the default logger instance to the given logger.
func setDefLogger(l *logger) {
	atomic.StorePointer(&def, unsafe.Pointer(l))
}

// GetLogger returns the default logger instance.
func GetLogger() Logger {
	return getDefLogger()
}

// SetLoggerOptions sets the options for the default logger.
func SetLoggerOptions(opts ...LoggerOption) {
	if len(opts) == 0 {
		return
	}
	l := getDefLogger().clone()
	for _, opt := range opts {
		opt(l)
	}
	setDefLogger(l)
}

// SetAppName sets the name of the application.
func SetAppName(name string) {
	l := getDefLogger().clone()
	l.app = EscapedString(name)
	setDefLogger(l)
}

// SetLevel sets the logging level for the default logger.
func SetLevel(level Level) {
	l := getDefLogger().clone()
	l.level = level
	setDefLogger(l)
}

// SetCaller sets whether or not to log the caller's function name and line number for the default logger.
func SetCaller(enable bool) {
	l := getDefLogger().clone()
	if enable {
		l.caller = Enable
	} else {
		l.caller = Disable
	}
	setDefLogger(l)
}

// SetColor sets whether or not to use colorized output levelTag on plain encoding for the default logger.
func SetColor(enable bool) {
	l := getDefLogger().clone()
	if enable {
		l.color = Enable
	} else {
		l.color = Disable
	}
	setDefLogger(l)
}

// SetShortFile sets whether or not to log the short file name for the default logger.
func SetShortFile(enable bool) {
	l := getDefLogger().clone()
	if enable {
		l.shortFile = Enable
	} else {
		l.shortFile = Disable
	}
	setDefLogger(l)
}

// SetTimeFormat sets the time format string for the default logger.
func SetTimeFormat(format string) {
	l := getDefLogger().clone()
	l.timeFmt = EscapedString(format)
	setDefLogger(l)
}

// SetEncode sets the log encoding type for the default logger.
func SetEncode(e EncodeType) {
	l := getDefLogger().clone()
	switch e {
	case PLAIN, JSON:
		l.encType = e
	default:
		l.encType = JSON
	}
	setDefLogger(l)
}

// SetEncodeFunc sets the log encoding type and encode function for the default logger.
func SetEncodeFunc(e EncodeFunc) {
	l := getDefLogger().clone()
	l.encType = -1
	l.enc = e
	setDefLogger(l)
}

// SetWriter sets the log writer for the default logger.
func SetWriter(w Writer) {
	l := getDefLogger().clone()
	l.wr = w
	setDefLogger(l)
}

// Log writes a log message with the given log level.
func Log(r Record) {
	l := getDefLogger()
	if l.IsEnabled(r.Level) {
		l.log(r)
	}
}

// Fatal logs a message at fatal level and exits the program with an error status.
func Fatal(args ...any) {
	getDefLogger().fatal(args...)
}

// Fatalf logs a formatted message at fatal level and exits the program with an error status.
func Fatalf(format string, args ...any) {
	getDefLogger().fatalf(format, args...)
}

// Fatalw logs a message with extra fields at fatal level and exits the program with an error status.
func Fatalw(msg string, fields ...Field) {
	getDefLogger().fatalw(msg, fields...)
}

// Error logs a message at error level.
func Error(args ...any) {
	getDefLogger().error(args...)
}

// Errorf logs a formatted message at error level.
func Errorf(format string, args ...any) {
	getDefLogger().errorf(format, args...)
}

// Errorw logs a message with extra fields at error level.
func Errorw(msg string, fields ...Field) {
	getDefLogger().errorw(msg, fields...)
}

// Warn logs a message at warning level.
func Warn(args ...any) {
	getDefLogger().warn(args...)
}

// Warnf logs a formatted message at warning level.
func Warnf(format string, args ...any) {
	getDefLogger().warnf(format, args...)
}

// Warnw logs a message with extra fields at warning level.
func Warnw(msg string, fields ...Field) {
	getDefLogger().warnw(msg, fields...)
}

// Notice logs a message at notice level.
func Notice(args ...any) {
	getDefLogger().notice(args...)
}

// Noticef logs a formatted message at notice level.
func Noticef(format string, args ...any) {
	getDefLogger().noticef(format, args...)
}

// Noticew logs a message with extra fields at notice level.
func Noticew(msg string, fields ...Field) {
	getDefLogger().noticew(msg, fields...)
}

// Info logs a message at info level.
func Info(args ...any) {
	getDefLogger().info(args...)
}

// Infof logs a formatted message at info level.
func Infof(format string, args ...any) {
	getDefLogger().infof(format, args...)
}

// Infow logs a message with extra fields at info level.
func Infow(msg string, fields ...Field) {
	getDefLogger().infow(msg, fields...)
}

// Debug logs a message at debug level.
func Debug(args ...any) {
	getDefLogger().debug(args...)
}

// Debugf logs a formatted message at debug level.
func Debugf(format string, args ...any) {
	getDefLogger().debugf(format, args...)
}

// Debugw logs a message with extra fields at debug level.
func Debugw(msg string, fields ...Field) {
	getDefLogger().debugw(msg, fields...)
}

// Trace logs a message at trace level.
func Trace(args ...any) {
	getDefLogger().trace(args...)
}

// Tracef logs a formatted message at trace level.
func Tracef(format string, args ...any) {
	getDefLogger().tracef(format, args...)
}

// Tracew logs a message with extra fields at trace level.
func Tracew(msg string, fields ...Field) {
	getDefLogger().tracew(msg, fields...)
}
