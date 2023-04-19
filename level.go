package olog

// Level is an unsigned 8-bit integer that represents the log level.
type Level uint8

// Define constants for each log level.
const (
	TRACE Level = iota
	DEBUG
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
)

// Define string representations of each log level.
const (
	tagTrace  = "trace"
	tagDebug  = "debug"
	tagInfo   = "info"
	tagNotice = "notice"
	tagWarn   = "warn"
	tagError  = "error"
	tagFatal  = "fatal"
)

// Define a mapping from Level values to their corresponding string representations.
var levelToStr = map[Level]string{
	TRACE:  tagTrace,
	DEBUG:  tagDebug,
	INFO:   tagInfo,
	NOTICE: tagNotice,
	WARN:   tagWarn,
	ERROR:  tagError,
	FATAL:  tagFatal,
}

// String returns the string representation of the Level value.
func (l Level) String() string {
	return levelToStr[l]
}
