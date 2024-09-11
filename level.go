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

// Define a mapping from string representations to their corresponding Level values.
var strToLevel = map[string]Level{
	tagTrace:  TRACE,
	tagDebug:  DEBUG,
	tagInfo:   INFO,
	tagNotice: NOTICE,
	tagWarn:   WARN,
	"warning": WARN,
	tagError:  ERROR,
	"err":     ERROR,
	tagFatal:  FATAL,
}

// String returns the string representation of the Level value.
func (l Level) String() string {
	return levelToStr[l]
}

// GetLevelByString returns the Level value corresponding to the given string.
func GetLevelByString(s string) Level {
	return strToLevel[s]
}
