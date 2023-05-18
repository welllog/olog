package olog

import (
	"runtime"
	"time"
)

// EncodeType is an enumeration type for different encoding types.
type EncodeType int8

const (
	// JSON represents the JSON encoding type.
	JSON EncodeType = iota
	// PLAIN represents the plain text encoding type.
	PLAIN
	// sep is the separator between fields in the log message.
	sep = '\t'
)

func FilterFields(fields []Field) []Field {
	n := len(fields)
	if n == 0 {
		return fields
	}

	set := make(map[string]struct{}, n)
	var remain int
	for idx, field := range fields {
		if !isSkipField(set, field.Key) {
			fields[remain], fields[idx] = fields[idx], fields[remain]
			remain++
		}
	}
	return fields[:remain]
}

func isSkipField(keysSet map[string]struct{}, key string) bool {
	_, ok := filterField[key]
	if ok {
		return true
	}

	_, ok = keysSet[key]
	if ok {
		return true
	}

	keysSet[key] = struct{}{}
	return false
}

var (
	fieldTime    = "@timestamp"
	fieldLevel   = "level"
	fieldContent = "content"
	fieldCaller  = "caller"
	fieldStack   = "stack"

	filterField = map[string]struct{}{
		fieldTime:    {},
		fieldLevel:   {},
		fieldContent: {},
		fieldCaller:  {},
		fieldStack:   {},
	}
)

// jsonEncode to encode a Record object as JSON and write it to the given Writer.
func jsonEncode(r Record, buf *Buffer) {
	_, _ = buf.WriteString(`{"@timestamp":"`)
	buf.WriteTime(time.Now(), r.TimeFmt)
	_, _ = buf.WriteString(`","level":"`)
	_, _ = buf.WriteString(r.LevelTag)

	if r.App != "" {
		_, _ = buf.WriteString(`","app":"`)
		_, _ = buf.WriteString(r.App)
	}

	var (
		get    bool
		more   bool
		frame  runtime.Frame
		frames *runtime.Frames
	)

	if r.Stack.IsOpen() {
		frames = getCallerFrames(r.CallerSkip, r.StackSize)
		frame, more = frames.Next()
		get = true
	}

	if r.Caller.IsOpen() {
		var (
			file string
			line int
		)
		if get {
			file, line = frame.File, frame.Line
		} else {
			file, line = getCaller(r.CallerSkip)
		}
		_, _ = buf.WriteString(`","caller":"`)
		_, _ = buf.WriteString(shortFile(file))
		_, _ = buf.WriteString(`:`)
		buf.WriteInt64(int64(line))
	}

	_, _ = buf.WriteString(`","content":`)

	if len(r.MsgArgs) > 0 {
		if r.MsgOrFormat != "" {
			buf.WriteQuoteSprintf(r.MsgOrFormat, r.MsgArgs...)
		} else {
			buf.WriteQuoteSprint(r.MsgArgs...)
		}
	} else {
		buf.WriteQuoteString(r.MsgOrFormat)
	}

	set := make(map[string]struct{}, len(r.Fields))
	// Loop over the fields of the Record object and write them to the buffer as JSON.
	for _, field := range r.Fields {
		if !isSkipField(set, field.Key) {
			// Write the field key and value to the buffer as JSON.
			_, _ = buf.WriteString(`,"`)
			_, _ = buf.WriteString(field.Key)
			_, _ = buf.WriteString(`":`)
			buf.WriteAny(field.Value, true)
		}
	}

	if r.Stack.IsOpen() {
		_, _ = buf.WriteString(`,"stack":"`)
		if frame.File != "" {
			for {
				_, _ = buf.WriteString(`\n`)
				_, _ = buf.WriteString(frame.Function)
				_, _ = buf.WriteString(`\n\t`)
				_, _ = buf.WriteString(frame.File)
				_ = buf.WriteByte(':')
				buf.WriteInt64(int64(frame.Line))

				if !more {
					break
				}
				frame, more = frames.Next()
			}
		}
		_ = buf.WriteByte('"')
	}

	// Write the closing curly brace and newline character to the buffer.
	_, _ = buf.WriteString("}\n")
}

// plainEncode to encode a Record object as plain text and write it to the given Writer.
func plainEncode(r Record, buf *Buffer, enableColor bool) {
	buf.WriteTime(time.Now(), r.TimeFmt)
	_ = buf.WriteByte(sep)
	if enableColor {
		writeLevelWithColor(r.Level, r.LevelTag, buf)
	} else {
		_, _ = buf.WriteString(r.LevelTag)
	}
	_ = buf.WriteByte(sep)

	if r.App != "" {
		_, _ = buf.WriteString(r.App)
		_ = buf.WriteByte(sep)
	}

	var (
		get    bool
		more   bool
		frame  runtime.Frame
		frames *runtime.Frames
	)

	if r.Stack.IsOpen() {
		frames = getCallerFrames(r.CallerSkip, r.StackSize)
		frame, more = frames.Next()
		get = true
	}

	if r.Caller.IsOpen() {
		var (
			file string
			line int
		)
		if get {
			file, line = frame.File, frame.Line
		} else {
			file, line = getCaller(r.CallerSkip)
		}
		file = shortFile(file)

		_, _ = buf.WriteString(file)
		_ = buf.WriteByte(':')
		buf.WriteInt64(int64(line))
		_ = buf.WriteByte(sep)
	}

	if len(r.MsgArgs) > 0 {
		if r.MsgOrFormat != "" {
			buf.WriteSprintf(r.MsgOrFormat, r.MsgArgs...)
		} else {
			buf.WriteSprint(r.MsgArgs...)
		}
	} else {
		_, _ = buf.WriteString(r.MsgOrFormat)
	}

	set := make(map[string]struct{}, len(r.Fields))
	// Loop over the fields of the Record object and write them to the buffer as plain text.
	for _, field := range r.Fields {
		if !isSkipField(set, field.Key) {
			_ = buf.WriteByte(sep)
			_, _ = buf.WriteString(field.Key)
			_ = buf.WriteByte('=')
			buf.WriteAny(field.Value, false)
		}
	}

	if r.Stack.IsOpen() {
		_ = buf.WriteByte(sep)
		_, _ = buf.WriteString("stack=")
		if frame.File != "" {
			for {
				_ = buf.WriteByte('\n')
				_, _ = buf.WriteString(frame.Function)
				_, _ = buf.WriteString("\n\t")
				_, _ = buf.WriteString(frame.File)
				_ = buf.WriteByte(':')
				buf.WriteInt64(int64(frame.Line))

				if !more {
					break
				}
				frame, more = frames.Next()
			}
		}
	}

	// Write the newline character to the buffer.
	_ = buf.WriteByte('\n')
}
