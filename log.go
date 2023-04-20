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

// SetLevel sets the logging level for the default logger.
func SetLevel(level Level) {
	l := getDefLogger().clone()
	l.level = level
	setDefLogger(l)
}

// SetCaller sets whether or not to log the caller's function name and line number for the default logger.
func SetCaller(enable bool) {
	l := getDefLogger().clone()
	l.enableCaller = enable
	setDefLogger(l)
}

// SetColor sets whether or not to use colorized output levelTag on plain encoding for the default logger.
func SetColor(enable bool) {
	l := getDefLogger().clone()
	l.enableColor = enable
	setDefLogger(l)
}

// SetTimeFormat sets the time format string for the default logger.
func SetTimeFormat(format string) {
	l := getDefLogger().clone()
	l.timeFormat = format
	setDefLogger(l)
}

// SetEncode sets the log encoding type for the default logger.
func SetEncode(e EncodeType) {
	l := getDefLogger().clone()
	l.encodeType = e
	setDefLogger(l)
}

// SetWriter sets the log writer for the default logger.
func SetWriter(w Writer) {
	l := getDefLogger().clone()
	l.writer = w
	setDefLogger(l)
}

func Log(opt LogOption, a ...any) {
	getDefLogger().log(opt, a...)
}

func Logf(opt LogOption, format string, a ...any) {
	getDefLogger().logf(opt, format, a...)
}

// Fatal logs a message at fatal level and exits the program with an error status.
func Fatal(a ...any) {
	getDefLogger().fatal(a...)
}

// Fatalf logs a formatted message at fatal level and exits the program with an error status.
func Fatalf(format string, a ...any) {
	getDefLogger().fatalf(format, a...)
}

// Fatalw logs a message with extra fields at fatal level and exits the program with an error status.
func Fatalw(msg string, fields ...Field) {
	getDefLogger().fatalw(msg, fields...)
}

// Error logs a message at error level.
func Error(a ...any) {
	getDefLogger().error(a...)
}

// Errorf logs a formatted message at error level.
func Errorf(format string, a ...any) {
	getDefLogger().errorf(format, a...)
}

// Errorw logs a message with extra fields at error level.
func Errorw(msg string, fields ...Field) {
	getDefLogger().errorw(msg, fields...)
}

// Warn logs a message at warning level.
func Warn(a ...any) {
	getDefLogger().warn(a...)
}

// Warnf logs a formatted message at warning level.
func Warnf(format string, a ...any) {
	getDefLogger().warnf(format, a...)
}

// Warnw logs a message with extra fields at warning level.
func Warnw(msg string, fields ...Field) {
	getDefLogger().warnw(msg, fields...)
}

// Notice logs a message at notice level.
func Notice(a ...any) {
	getDefLogger().notice(a...)
}

// Noticef logs a formatted message at notice level.
func Noticef(format string, a ...any) {
	getDefLogger().noticef(format, a...)
}

// Noticew logs a message with extra fields at notice level.
func Noticew(msg string, fields ...Field) {
	getDefLogger().noticew(msg, fields...)
}

// Info logs a message at info level.
func Info(a ...any) {
	getDefLogger().info(a...)
}

// Infof logs a formatted message at info level.
func Infof(format string, a ...any) {
	getDefLogger().infof(format, a...)
}

// Infow logs a message with extra fields at info level.
func Infow(msg string, fields ...Field) {
	getDefLogger().infow(msg, fields...)
}

// Debug logs a message at debug level.
func Debug(a ...any) {
	getDefLogger().debug(a...)
}

// Debugf logs a formatted message at debug level.
func Debugf(format string, a ...any) {
	getDefLogger().debugf(format, a...)
}

// Debugw logs a message with extra fields at debug level.
func Debugw(msg string, fields ...Field) {
	getDefLogger().debugw(msg, fields...)
}

// Trace logs a message at trace level.
func Trace(a ...any) {
	getDefLogger().trace(a...)
}

// Tracef logs a formatted message at trace level.
func Tracef(format string, a ...any) {
	getDefLogger().tracef(format, a...)
}

// Tracew logs a message with extra fields at trace level.
func Tracew(msg string, fields ...Field) {
	getDefLogger().tracew(msg, fields...)
}
