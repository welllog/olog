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

type stringWriter interface {
	WriteString(string) (int, error)
}

// writeLevelWithColor takes in a level of logging and a tag string, and returns a string that
// contains the tag string wrapped with an ANSI color code to represent the level of logging.
// The returned string will have different colors depending on the level of logging.
func writeLevelWithColor(level Level, tag string, w stringWriter) {
	switch level {
	case FATAL:
		_, _ = w.WriteString(RedBold)
	case ERROR:
		_, _ = w.WriteString(Red)
	case WARN:
		_, _ = w.WriteString(Yellow)
	case NOTICE:
		_, _ = w.WriteString(Blue)
	case INFO:
		_, _ = w.WriteString(Green)
	case DEBUG:
		_, _ = w.WriteString(Gray)
	case TRACE:
		_, _ = w.WriteString(Cyan)
	default:
		_, _ = w.WriteString(tag)
		return
	}
	_, _ = w.WriteString(tag)
	_, _ = w.WriteString(Reset)
}
