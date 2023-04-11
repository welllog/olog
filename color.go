//go:build !windows

package olog

const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	RedBold   = "\033[1;31m"
	Green     = "\033[32m"
	GreenBold = "\033[1;32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	BlueBold  = "\033[1;34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
)

// writeLevelWithColor takes in a level of logging and a tag string, and returns a string that
// contains the tag string wrapped with an ANSI color code to represent the level of logging.
// The returned string will have different colors depending on the level of logging.
func writeLevelWithColor(level Level, tag string, buf *Buffer) {
	switch level {
	case FATAL:
		_, _ = buf.WriteString(RedBold)
	case ERROR:
		_, _ = buf.WriteString(Red)
	case WARN:
		_, _ = buf.WriteString(Yellow)
	case INFO:
		_, _ = buf.WriteString(Green)
	case DEBUG:
		_, _ = buf.WriteString(Gray)
	default:
		_, _ = buf.WriteString(tag)
		return
	}
	_, _ = buf.WriteString(tag)
	_, _ = buf.WriteString(Reset)
}
