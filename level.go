package olog

import "strings"

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

// Use an array instead of a map for faster lookup.
var levelTags = [...]string{
	tagTrace,
	tagDebug,
	tagInfo,
	tagNotice,
	tagWarn,
	tagError,
	tagFatal,
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
	if l <= FATAL {
		return levelTags[l]
	}
	return ""
}

// GetLevelByString returns the Level value corresponding to the given string.
func GetLevelByString(s string) Level {
	return strToLevel[strings.ToLower(s)]
}
