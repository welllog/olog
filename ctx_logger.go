package olog

import (
	"context"
	"fmt"
	"sort"
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
	if len(fields) == 0 {
		return logger
	}
	fields = uniqueAscFields(fields)
	return withFields(logger, fields...)
}

// WithEntries creates a new logger with the provided entries.
func WithEntries(logger Logger, entries map[string]any) Logger {
	if len(entries) == 0 {
		return logger
	}

	fields := make([]Field, 0, len(entries))
	for k, v := range entries {
		fields = append(fields, Field{Key: k, Value: v})
	}
	sort.Sort(AscFields(fields))
	return withFields(logger, fields...)
}

// WithFields creates a new logger with the provided fields.
// It will update fields
func WithFields(logger Logger, fields ...Field) Logger {
	if len(fields) == 0 {
		return logger
	}

	fields = uniqueAscFields(fields)
	return withFields(logger, fields...)
}

// WithKV creates a new logger with the provided key-value pairs.
// If the number of key-value pairs is odd, the last key will have the value "KEYVALS UNPAIRED".
func WithKV(logger Logger, kv ...any) Logger {
	if len(kv) == 0 {
		return logger
	}

	n := len(kv)
	if (n & 1) == 1 {
		n++
	}
	fields := make([]Field, 0, n>>1)

	if n == len(kv) {
		for i := 0; i < len(kv); i += 2 {
			fields = append(fields, Field{Key: toString(kv[i]), Value: kv[i+1]})
		}
	} else {
		for i := 0; i < len(kv)-1; i += 2 {
			fields = append(fields, Field{Key: toString(kv[i]), Value: kv[i+1]})
		}
		fields = append(fields, Field{Key: toString(kv[len(kv)-1]), Value: "KEYVALS UNPAIRED"})
	}
	fields = uniqueAscFields(fields)
	return withFields(logger, fields...)
}

func withFields(logger Logger, fields ...Field) Logger {
	innerCtx, ok := logger.(*ctxLogger)
	if !ok {
		return &ctxLogger{
			Logger: logger,
			fields: fields,
		}
	}

	innerFields := innerCtx.fields
	if len(innerFields) == 0 {
		return &ctxLogger{
			Logger: innerCtx.Logger,
			fields: fields,
		}
	}

	// Calculate the number of fields in innerFields that are not in fields.
	var n int
	i, j := 0, 0
	for i < len(innerFields) && j < len(fields) {
		if innerFields[i].Key == fields[j].Key {
			i++
			j++
		} else if innerFields[i].Key < fields[j].Key {
			n++
			i++
		} else {
			j++
		}
	}
	n += len(innerFields) - i

	if n == 0 {
		return &ctxLogger{
			Logger: innerCtx.Logger,
			fields: fields,
		}
	}

	merged := make([]Field, 0, len(fields)+n)
	i, j = 0, 0
	for i < len(innerFields) && j < len(fields) {
		if innerFields[i].Key == fields[j].Key {
			merged = append(merged, fields[j])
			i++
			j++
		} else if innerFields[i].Key < fields[j].Key {
			merged = append(merged, innerFields[i])
			i++
		} else {
			merged = append(merged, fields[j])
			j++
		}
	}
	merged = append(merged, innerFields[i:]...)
	merged = append(merged, fields[j:]...)

	return &ctxLogger{
		Logger: innerCtx.Logger,
		fields: merged,
	}
}

func (c *ctxLogger) Log(r Record) {
	if c.IsEnabled(r.Level) {
		r.Fields = c.buildFields(r.Fields...)
		r.CallerSkip++
		c.Logger.Log(r)
	}
}

func (c *ctxLogger) Fatal(args ...any) {
	c.Logger.Log(Record{
		Level:      FATAL,
		CallerSkip: 1,
		MsgArgs:    args,
		Fields:     c.fields,
		OsExit:     true,
	})
}

func (c *ctxLogger) Fatalf(format string, args ...any) {
	c.Logger.Log(Record{
		Level:       FATAL,
		CallerSkip:  1,
		MsgOrFormat: format,
		MsgArgs:     args,
		Fields:      c.fields,
		OsExit:      true,
	})
}

func (c *ctxLogger) Fatalw(msg string, fields ...Field) {
	c.Logger.Log(Record{
		Level:       FATAL,
		CallerSkip:  1,
		MsgOrFormat: msg,
		Fields:      c.buildFields(fields...),
		OsExit:      true,
	})
}

func (c *ctxLogger) Error(args ...any) {
	if c.IsEnabled(ERROR) {
		c.Logger.Log(Record{
			Level:      ERROR,
			CallerSkip: 1,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Errorf(format string, args ...any) {
	if c.IsEnabled(ERROR) {
		c.Logger.Log(Record{
			Level:       ERROR,
			CallerSkip:  1,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Errorw(msg string, fields ...Field) {
	if c.IsEnabled(ERROR) {
		c.Logger.Log(Record{
			Level:       ERROR,
			CallerSkip:  1,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Warn(args ...any) {
	if c.IsEnabled(WARN) {
		c.Logger.Log(Record{
			Level:      WARN,
			CallerSkip: 1,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Warnf(format string, args ...any) {
	if c.IsEnabled(WARN) {
		c.Logger.Log(Record{
			Level:       WARN,
			CallerSkip:  1,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Warnw(msg string, fields ...Field) {
	if c.IsEnabled(WARN) {
		c.Logger.Log(Record{
			Level:       WARN,
			CallerSkip:  1,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Notice(args ...any) {
	if c.IsEnabled(NOTICE) {
		c.Logger.Log(Record{
			Level:      NOTICE,
			CallerSkip: 1,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Noticef(format string, args ...any) {
	if c.IsEnabled(NOTICE) {
		c.Logger.Log(Record{
			Level:       NOTICE,
			CallerSkip:  1,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Noticew(msg string, fields ...Field) {
	if c.IsEnabled(NOTICE) {
		c.Logger.Log(Record{
			Level:       NOTICE,
			CallerSkip:  1,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Info(args ...any) {
	if c.IsEnabled(INFO) {
		c.Logger.Log(Record{
			Level:      INFO,
			CallerSkip: 1,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Infof(format string, args ...any) {
	if c.IsEnabled(INFO) {
		c.Logger.Log(Record{
			Level:       INFO,
			CallerSkip:  1,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Infow(msg string, fields ...Field) {
	if c.IsEnabled(INFO) {
		c.Logger.Log(Record{
			Level:       INFO,
			CallerSkip:  1,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Debug(args ...any) {
	if c.IsEnabled(DEBUG) {
		c.Logger.Log(Record{
			Level:      DEBUG,
			CallerSkip: 1,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Debugf(format string, args ...any) {
	if c.IsEnabled(DEBUG) {
		c.Logger.Log(Record{
			Level:       DEBUG,
			CallerSkip:  1,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Debugw(msg string, fields ...Field) {
	if c.IsEnabled(DEBUG) {
		c.Logger.Log(Record{
			Level:       DEBUG,
			CallerSkip:  1,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

func (c *ctxLogger) Trace(args ...any) {
	if c.IsEnabled(TRACE) {
		c.Logger.Log(Record{
			Level:      TRACE,
			CallerSkip: 1,
			Stack:      Enable,
			MsgArgs:    args,
			Fields:     c.fields,
		})
	}
}

func (c *ctxLogger) Tracef(format string, args ...any) {
	if c.IsEnabled(TRACE) {
		c.Logger.Log(Record{
			Level:       TRACE,
			CallerSkip:  1,
			Stack:       Enable,
			MsgOrFormat: format,
			MsgArgs:     args,
			Fields:      c.fields,
		})
	}
}

func (c *ctxLogger) Tracew(msg string, fields ...Field) {
	if c.IsEnabled(TRACE) {
		c.Logger.Log(Record{
			Level:       TRACE,
			CallerSkip:  1,
			Stack:       Enable,
			MsgOrFormat: msg,
			Fields:      c.buildFields(fields...),
		})
	}
}

// buildFields builds the final fields slice.
func (c *ctxLogger) buildFields(fields ...Field) []Field {
	// No new fields, reuse old fields
	if len(fields) == 0 {
		return c.fields
	}

	// not modify the old fields, and ensure that the new field is in front.
	return append(fields, c.fields...)
}

func uniqueAscFields(fields []Field) []Field {
	sort.Sort(AscFields(fields))
	dst := fields[:1]
	for i := 1; i < len(fields); i++ {
		if fields[i].Key != fields[i-1].Key {
			dst = append(dst, fields[i])
		}
	}
	return dst
}

func toString(key any) string {
	switch k := key.(type) {
	case string:
		return k
	case fmt.Stringer:
		return k.String()
	default:
		return fmt.Sprint(k)
	}
}
