package olog

import (
	"context"
)

// The ctxLogger is a logger that can be used to log with context.
type ctxLogger struct {
	Logger
	fields []Field
}

// WithContext creates a new logger with the provided context.
// If handles is not provided, the default context handle is used.
func WithContext(logger Logger, ctx context.Context, handles ...CtxHandle) Logger {
	var handle CtxHandle
	if len(handles) > 0 {
		handle = handles[0]
	} else {
		handle = getDefCtxHandle()
	}

	fields := handle(ctx)

	return &ctxLogger{
		Logger: logger,
		fields: logger.buildFields(fields...),
	}
}

// WithEntries creates a new logger with the provided entries.
func WithEntries(logger Logger, entries map[string]interface{}) Logger {
	fields := make([]Field, 0, len(entries))
	for k, v := range entries {
		fields = append(fields, Field{Key: k, Value: v})
	}

	return &ctxLogger{
		Logger: logger,
		fields: logger.buildFields(fields...),
	}
}

func (c *ctxLogger) Log(r Record) {
	if c.IsEnabled(r.Level) {
		r.Fields = c.buildFields(r.Fields...)
		c.log(r)
	}
}

func (c *ctxLogger) Fatal(a ...interface{}) {
	c.log(Record{
		Level:   FATAL,
		MsgArgs: a,
		Fields:  c.fields,
		OsExit:  true,
	})
}

func (c *ctxLogger) Fatalf(format string, a ...interface{}) {
	c.log(Record{
		Level:       FATAL,
		MsgOrFormat: format,
		MsgArgs:     a,
		Fields:      c.fields,
		OsExit:      true,
	})
}

func (c *ctxLogger) Fatalw(msg string, fields ...Field) {
	c.log(Record{
		Level:       FATAL,
		MsgOrFormat: msg,
		Fields:      c.buildFields(fields...),
		OsExit:      true,
	})
}

func (c *ctxLogger) Error(a ...interface{}) {
	if c.IsEnabled(ERROR) {
		c.log(Record{
			Level:   ERROR,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Errorf(format string, a ...interface{}) {
	if c.IsEnabled(ERROR) {
		c.log(Record{
			Level:       ERROR,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

func (c *ctxLogger) Warn(a ...interface{}) {
	if c.IsEnabled(WARN) {
		c.log(Record{
			Level:   WARN,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Warnf(format string, a ...interface{}) {
	if c.IsEnabled(WARN) {
		c.log(Record{
			Level:       WARN,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

func (c *ctxLogger) Notice(a ...interface{}) {
	if c.IsEnabled(NOTICE) {
		c.log(Record{
			Level:   NOTICE,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Noticef(format string, a ...interface{}) {
	if c.IsEnabled(NOTICE) {
		c.log(Record{
			Level:       NOTICE,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

func (c *ctxLogger) Info(a ...interface{}) {
	if c.IsEnabled(INFO) {
		c.log(Record{
			Level:   INFO,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Infof(format string, a ...interface{}) {
	if c.IsEnabled(INFO) {
		c.log(Record{
			Level:       INFO,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

func (c *ctxLogger) Debug(a ...interface{}) {
	if c.IsEnabled(DEBUG) {
		c.log(Record{
			Level:   DEBUG,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Debugf(format string, a ...interface{}) {
	if c.IsEnabled(DEBUG) {
		c.log(Record{
			Level:       DEBUG,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

func (c *ctxLogger) Trace(a ...interface{}) {
	if c.IsEnabled(TRACE) {
		c.log(Record{
			Level:   TRACE,
			Stack:   Enable,
			MsgArgs: a,
			Fields:  c.fields,
		})
	}
}

func (c *ctxLogger) Tracef(format string, a ...interface{}) {
	if c.IsEnabled(TRACE) {
		c.log(Record{
			Level:       TRACE,
			Stack:       Enable,
			MsgOrFormat: format,
			MsgArgs:     a,
			Fields:      c.fields,
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

// buildFields builds the final fields slice.
func (c *ctxLogger) buildFields(fields ...Field) []Field {
	if len(fields) == 0 {
		return c.fields
	}

	return append(fields, c.fields...)
}
