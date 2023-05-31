//go:build windows
// +build windows

package olog

type stringWriter interface {
	WriteString(string) (int, error)
}

// writeLevelWithColor is a stub function that always returns the input tag string unchanged. It's used as a placeholder
// function on Windows OS because the ANSI color codes that are used to format console output on Unix-based systems aren't
// supported by the Windows console.
func writeLevelWithColor(level Level, tag string, w stringWriter) {
	_, _ = w.WriteString(tag)
}
