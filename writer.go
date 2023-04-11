package olog

import (
	"io"
	"os"
)

// csWriter is a variable that holds a new instance of consoleWriter created by calling the NewConsoleWriter function
var csWriter = NewConsoleWriter()

// Writer is an interface that defines a Write method with a Level and a byte slice as its parameters and returns an integer and an error
// Special attention must be paid to the fact that p []byte should not exceed the scope of the Write method.
// After the Write method ends, the byte slice should not be used, otherwise will cause memory data errors.
type Writer interface {
	Write(level Level, p []byte) (n int, err error)
}

// consoleWriter is a struct that holds a standard output writer and a standard error writer
type consoleWriter struct {
	sw io.Writer
	ew io.Writer
}

// Write is a method on consoleWriter that writes the byte slice p to the standard writer or the error writer depending on the level parameter
func (c *consoleWriter) Write(level Level, p []byte) (n int, err error) {
	if level >= WARN {
		return c.ew.Write(p)
	}
	return c.sw.Write(p)
}

// customWriter is a struct that holds a custom writer
type customWriter struct {
	w io.Writer
}

// Write is a method on customWriter that writes the byte slice p to the custom writer
func (c *customWriter) Write(level Level, p []byte) (n int, err error) {
	return c.w.Write(p)
}

// NewConsoleWriter is a function that creates a new consoleWriter with os.Stdout as the standard writer and os.Stderr as the error writer
func NewConsoleWriter() Writer {
	return &consoleWriter{
		sw: os.Stdout,
		ew: os.Stderr,
	}
}

// NewWriter is a function that creates a new customWriter with the specified writer as its parameter
// Special attention must be paid to the fact that []byte should not exceed the scope of the Write method.
// After the Write method ends, the byte slice should not be used, otherwise will cause memory data errors.
func NewWriter(w io.Writer) Writer {
	return &customWriter{
		w: w,
	}
}
