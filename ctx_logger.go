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

func (c *ctxLogger) Log(r Record) {
	if c.IsEnabled(r.Level) {
		r.Fields = c.buildFields(r.Fields...)
		c.log(r)
	}
}

func (c *ctxLogger) Fatal(a ...any) {
	if c.IsEnabled(FATAL) {
		c.log(Record{
			Level:   FATAL,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Fatalf(format string, a ...any) {
	if c.IsEnabled(FATAL) {
		c.log(Record{
			Level:       FATAL,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Fatalw(msg string, fields ...Field) {
	if c.IsEnabled(FATAL) {
		c.log(Record{
			Level:       FATAL,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Error(a ...any) {
	if c.IsEnabled(ERROR) {
		c.log(Record{
			Level:   ERROR,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Errorf(format string, a ...any) {
	if c.IsEnabled(ERROR) {
		c.log(Record{
			Level:       ERROR,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Errorw(msg string, fields ...Field) {
	if c.IsEnabled(ERROR) {
		c.log(Record{
			Level:       ERROR,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Warn(a ...any) {
	if c.IsEnabled(WARN) {
		c.log(Record{
			Level:   WARN,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Warnf(format string, a ...any) {
	if c.IsEnabled(WARN) {
		c.log(Record{
			Level:       WARN,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Warnw(msg string, fields ...Field) {
	if c.IsEnabled(WARN) {
		c.log(Record{
			Level:       WARN,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Notice(a ...any) {
	if c.IsEnabled(NOTICE) {
		c.log(Record{
			Level:   NOTICE,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Noticef(format string, a ...any) {
	if c.IsEnabled(NOTICE) {
		c.log(Record{
			Level:       NOTICE,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Noticew(msg string, fields ...Field) {
	if c.IsEnabled(NOTICE) {
		c.log(Record{
			Level:       NOTICE,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Info(a ...any) {
	if c.IsEnabled(INFO) {
		c.log(Record{
			Level:   INFO,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Infof(format string, a ...any) {
	if c.IsEnabled(INFO) {
		c.log(Record{
			Level:       INFO,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Infow(msg string, fields ...Field) {
	if c.IsEnabled(INFO) {
		c.log(Record{
			Level:       INFO,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Debug(a ...any) {
	if c.IsEnabled(DEBUG) {
		c.log(Record{
			Level:   DEBUG,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Debugf(format string, a ...any) {
	if c.IsEnabled(DEBUG) {
		c.log(Record{
			Level:       DEBUG,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Debugw(msg string, fields ...Field) {
	if c.IsEnabled(DEBUG) {
		c.log(Record{
			Level:       DEBUG,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Trace(a ...any) {
	if c.IsEnabled(TRACE) {
		c.log(Record{
			Level:   TRACE,
			Stack:   Enable,
			MsgArgs: a,
			Fields:  c.buildFields(),
		})
	}
}

func (c *ctxLogger) Tracef(format string, a ...any) {
	if c.IsEnabled(TRACE) {
		c.log(Record{
			Level:       TRACE,
			Stack:       Enable,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.buildFields(),
		})
	}
}

func (c *ctxLogger) Tracew(msg string, fields ...Field) {
	if c.IsEnabled(TRACE) {
		c.log(Record{
			Level:       TRACE,
			Stack:       Enable,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
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
