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

func (c *ctxLogger) Log(opt LogOption, a ...any) {
	if c.IsEnabled(opt.Level) {
		opt.Fields = c.buildFields(opt.Fields...)
		c.log(opt, a...)
	}
}

func (c *ctxLogger) Logf(opt LogOption, format string, a ...any) {
	if c.IsEnabled(opt.Level) {
		opt.Fields = c.buildFields(opt.Fields...)
		c.logf(opt, format, a...)
	}
}

func (c *ctxLogger) Fatal(a ...any) {
	if c.IsEnabled(FATAL) {
		c.log(LogOption{
			Level:  FATAL,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Fatalf(format string, a ...any) {
	if c.IsEnabled(FATAL) {
		c.logf(LogOption{
			Level:  FATAL,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Fatalw(msg string, fields ...Field) {
	if c.IsEnabled(FATAL) {
		c.logf(LogOption{
			Level:  FATAL,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Error(a ...any) {
	if c.IsEnabled(ERROR) {
		c.log(LogOption{
			Level:  ERROR,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Errorf(format string, a ...any) {
	if c.IsEnabled(ERROR) {
		c.logf(LogOption{
			Level:  ERROR,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Errorw(msg string, fields ...Field) {
	if c.IsEnabled(ERROR) {
		c.logf(LogOption{
			Level:  ERROR,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Warn(a ...any) {
	if c.IsEnabled(WARN) {
		c.log(LogOption{
			Level:  WARN,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Warnf(format string, a ...any) {
	if c.IsEnabled(WARN) {
		c.logf(LogOption{
			Level:  WARN,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Warnw(msg string, fields ...Field) {
	if c.IsEnabled(WARN) {
		c.logf(LogOption{
			Level:  WARN,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Notice(a ...any) {
	if c.IsEnabled(NOTICE) {
		c.log(LogOption{
			Level:  NOTICE,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Noticef(format string, a ...any) {
	if c.IsEnabled(NOTICE) {
		c.logf(LogOption{
			Level:  NOTICE,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Noticew(msg string, fields ...Field) {
	if c.IsEnabled(NOTICE) {
		c.logf(LogOption{
			Level:  NOTICE,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Info(a ...any) {
	if c.IsEnabled(INFO) {
		c.log(LogOption{
			Level:  INFO,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Infof(format string, a ...any) {
	if c.IsEnabled(INFO) {
		c.logf(LogOption{
			Level:  INFO,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Infow(msg string, fields ...Field) {
	if c.IsEnabled(INFO) {
		c.logf(LogOption{
			Level:  INFO,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Debug(a ...any) {
	if c.IsEnabled(DEBUG) {
		c.log(LogOption{
			Level:  DEBUG,
			Fields: c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Debugf(format string, a ...any) {
	if c.IsEnabled(DEBUG) {
		c.logf(LogOption{
			Level:  DEBUG,
			Fields: c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Debugw(msg string, fields ...Field) {
	if c.IsEnabled(DEBUG) {
		c.logf(LogOption{
			Level:  DEBUG,
			Fields: c.buildFields(fields...),
		}, msg)
	}
}

func (c *ctxLogger) Trace(a ...any) {
	if c.IsEnabled(TRACE) {
		c.log(LogOption{
			Level:       TRACE,
			EnableStack: EnableOpen,
			Fields:      c.buildFields(),
		}, a...)
	}
}

func (c *ctxLogger) Tracef(format string, a ...any) {
	if c.IsEnabled(TRACE) {
		c.logf(LogOption{
			Level:       TRACE,
			EnableStack: EnableOpen,
			Fields:      c.buildFields(),
		}, format, a...)
	}
}

func (c *ctxLogger) Tracew(msg string, fields ...Field) {
	if c.IsEnabled(TRACE) {
		c.logf(LogOption{
			Level:       TRACE,
			EnableStack: EnableOpen,
			Fields:      c.buildFields(fields...),
		}, msg)
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
