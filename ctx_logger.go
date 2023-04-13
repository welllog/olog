package olog

import (
	"context"
)

// The ctxLogger type embeds the Logger interface and adds context, entries, and context handle fields.
type ctxLogger struct {
	Logger
	ctx       context.Context
	entries   map[string]any
	ctxHandle CtxHandle
}

// WithContext creates a new logger with the provided context and context handle.
// If handles is not provided, the default context handle is used.
func WithContext(logger Logger, ctx context.Context, handles ...CtxHandle) Logger {
	var handle CtxHandle
	if len(handles) > 0 {
		handle = handles[0]
	} else {
		handle = getDefCtxHandle()
	}
	return &ctxLogger{
		Logger:    logger,
		ctx:       ctx,
		ctxHandle: handle,
	}
}

// WithEntries creates a new logger with the provided entries and an empty context.
func WithEntries(logger Logger, entries map[string]any) Logger {
	return &ctxLogger{
		Logger:    logger,
		ctx:       context.Background(),
		entries:   entries,
		ctxHandle: emptyHandle,
	}
}

func (c *ctxLogger) Log(level Level, opts ...LogOption) {
	if c.IsEnabled(level) {
		opts = append(opts, WithFields(c.buildFields()...))
		c.log(level, opts...)
	}
}

func (c *ctxLogger) Fatal(a ...any) {
	if c.IsEnabled(FATAL) {
		c.log(FATAL, WithPrint(a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Fatalf(format string, a ...any) {
	if c.IsEnabled(FATAL) {
		c.log(FATAL, WithPrintf(format, a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Fatalw(msg string, fields ...Field) {
	if c.IsEnabled(FATAL) {
		c.log(FATAL, WithPrintMsg(msg), WithFields(c.buildFields(fields...)...))
	}
}

func (c *ctxLogger) Error(a ...any) {
	if c.IsEnabled(ERROR) {
		c.log(ERROR, WithPrint(a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Errorf(format string, a ...any) {
	if c.IsEnabled(ERROR) {
		c.log(ERROR, WithPrintf(format, a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Errorw(msg string, fields ...Field) {
	if c.IsEnabled(ERROR) {
		c.log(ERROR, WithPrintMsg(msg), WithFields(c.buildFields(fields...)...))
	}
}

func (c *ctxLogger) Warn(a ...any) {
	if c.IsEnabled(WARN) {
		c.log(WARN, WithPrint(a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Warnf(format string, a ...any) {
	if c.IsEnabled(WARN) {
		c.log(WARN, WithPrintf(format, a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Warnw(msg string, fields ...Field) {
	if c.IsEnabled(WARN) {
		c.log(WARN, WithPrintMsg(msg), WithFields(c.buildFields(fields...)...))
	}
}

func (c *ctxLogger) Info(a ...any) {
	if c.IsEnabled(INFO) {
		c.log(INFO, WithPrint(a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Infof(format string, a ...any) {
	if c.IsEnabled(INFO) {
		c.log(INFO, WithPrintf(format, a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Infow(msg string, fields ...Field) {
	if c.IsEnabled(INFO) {
		c.log(INFO, WithPrintMsg(msg), WithFields(c.buildFields(fields...)...))
	}
}

func (c *ctxLogger) Debug(a ...any) {
	if c.IsEnabled(DEBUG) {
		c.log(DEBUG, WithPrint(a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Debugf(format string, a ...any) {
	if c.IsEnabled(DEBUG) {
		c.log(DEBUG, WithPrintf(format, a...), WithFields(c.buildFields()...))
	}
}

func (c *ctxLogger) Debugw(msg string, fields ...Field) {
	if c.IsEnabled(DEBUG) {
		c.log(DEBUG, WithPrintMsg(msg), WithFields(c.buildFields(fields...)...))
	}
}

// buildFields appends fields from the logger's context and entries, then returns the resulting fields slice.
func (c *ctxLogger) buildFields(fields ...Field) []Field {
	// Retrieve context fields and compute minimum capacity needed to hold all fields.
	ctxFields := c.ctxHandle(c.ctx)
	minCap := len(fields) + len(c.entries) + len(ctxFields)

	// Check if fields slice capacity is sufficient for holding all fields.
	if cap(fields) < minCap {
		// If not, allocate a new slice with double the capacity or the minimum capacity, whichever is greater.
		newCap := 2 * cap(fields)
		if newCap < minCap {
			newCap = minCap
		}
		newFields := make([]Field, len(fields), newCap)
		// Copy existing fields to the new slice and update the fields reference.
		copy(newFields, fields)
		fields = newFields
	}

	// Append context fields and entries to fields slice.
	fields = append(fields, ctxFields...)
	for key, entry := range c.entries {
		fields = append(fields, Field{Key: key, Value: entry})
	}

	return c.Logger.buildFields(fields...)
}
