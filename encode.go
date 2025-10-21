package olog

import (
	"runtime"

	"github.com/welllog/olog/encoder"
)

// EncodeType is an enumeration type for different encoding types.
type EncodeType int8

const (
	// JSON represents the JSON encoding type.
	JSON EncodeType = iota
	// PLAIN represents the plain text encoding type.
	PLAIN
)

func FilterFields(fields []Field) []Field {
	n := len(fields)
	if n == 0 {
		return fields
	}

	if n <= 4 {
		remain := 0
	outer:
		for i := 0; i < n; i++ {
			f := fields[i]
			if _, skip := filterField[f.Key]; skip {
				continue
			}
			for j := 0; j < remain; j++ {
				if fields[j].Key == f.Key {
					continue outer
				}
			}
			fields[remain], fields[i] = fields[i], fields[remain]
			remain++
		}
		return fields[:remain]
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
	fieldApp     = "app"
	fieldContent = "content"
	fieldCaller  = "caller"
	fieldStack   = "stack"

	filterField = map[string]struct{}{
		fieldTime:    {},
		fieldLevel:   {},
		fieldApp:     {},
		fieldContent: {},
		fieldCaller:  {},
		fieldStack:   {},
	}
)

// jsonEncode to encode a Record object as JSON to the buffer.
func jsonEncode(r Record, buf *encoder.Buffer) {
	enc := encoder.JsonEncoder{Buffer: buf}

	_, _ = enc.WriteString(`{"@timestamp":"`)
	enc.WriteTime(r.Time, r.TimeFmt)
	_, _ = enc.WriteString(`","level":"`)
	_, _ = enc.WriteString(r.LevelTag)

	if r.App != "" {
		_, _ = enc.WriteString(`","app":"`)
		_, _ = enc.WriteString(r.App)
	}

	var (
		more  bool
		frame runtime.Frame
	)
	frames := r.Frames()
	if frames != nil {
		frame, more = frames.Next()
	}

	if r.Caller.IsOpen() {
		file := frame.File
		if r.ShortFile.IsOpen() {
			file = shortFile(file)
		}
		_, _ = enc.WriteString(`","caller":"`)
		enc.WriteEscapedString(file)
		_, _ = enc.WriteString(`:`)
		enc.WriteInt64(int64(frame.Line))
	}

	_, _ = enc.WriteString(`","content":"`)
	_, _ = encoder.EPrintf(enc, r.MsgOrFormat, r.MsgArgs...)
	enc.WriteQuote()

outer:
	for i, field := range r.Fields {
		if _, ok := filterField[field.Key]; ok {
			continue
		}

		for j := 0; j < i; j++ {
			if field.Key == r.Fields[j].Key {
				continue outer
			}
		}

		_, _ = enc.WriteString(`,"`)
		enc.WriteEscapedString(field.Key)
		_, _ = enc.WriteString(`":`)
		enc.WriteValue(field.Value)
	}

	if r.Stack.IsOpen() {
		_, _ = enc.WriteString(`,"stack":"`)
		if frame.PC != 0 {
			for {
				_, _ = enc.WriteString(`\n`)
				enc.WriteEscapedString(frame.Function)
				_, _ = enc.WriteString(`\n\t`)
				enc.WriteEscapedString(frame.File)
				_ = enc.WriteByte(':')
				enc.WriteInt64(int64(frame.Line))

				if !more {
					break
				}
				frame, more = frames.Next()
			}
		}
		enc.WriteQuote()
	}

	// Write the closing curly brace and newline character to the buffer.
	_, _ = enc.WriteString("}\n")
}

// plainEncode to encode a Record object as plain text to the buffer.
func plainEncode(r Record, buf *encoder.Buffer, enableColor bool) {
	enc := encoder.PlainEncoder{Buffer: buf}

	enc.WriteTime(r.Time, r.TimeFmt)
	enc.WriteSeparator()
	if enableColor {
		writeLevelWithColor(r.Level, r.LevelTag, enc)
	} else {
		_, _ = enc.WriteString(r.LevelTag)
	}
	enc.WriteSeparator()

	if r.App != "" {
		_, _ = enc.WriteString(r.App)
		enc.WriteSeparator()
	}

	var (
		more  bool
		frame runtime.Frame
	)
	frames := r.Frames()
	if frames != nil {
		frame, more = frames.Next()
	}

	if r.Caller.IsOpen() {
		file := frame.File
		if r.ShortFile.IsOpen() {
			file = shortFile(file)
		}

		_, _ = enc.WriteString(file)
		_ = enc.WriteByte(':')
		enc.WriteInt64(int64(frame.Line))
		enc.WriteSeparator()
	}

	_, _ = encoder.EPrintf(enc, r.MsgOrFormat, r.MsgArgs...)

outer:
	for i, field := range r.Fields {
		if _, ok := filterField[field.Key]; ok {
			continue
		}

		for j := 0; j < i; j++ {
			if field.Key == r.Fields[j].Key {
				continue outer
			}
		}

		enc.WriteSeparator()
		enc.WriteName(field.Key)
		enc.WriteValue(field.Value)
	}

	if r.Stack.IsOpen() {
		enc.WriteSeparator()
		_, _ = enc.WriteString("stack=")
		if frame.PC != 0 {
			for {
				_ = enc.WriteByte('\n')
				_, _ = enc.WriteString(frame.Function)
				_, _ = enc.WriteString("\n\t")
				_, _ = enc.WriteString(frame.File)
				_ = enc.WriteByte(':')
				enc.WriteInt64(int64(frame.Line))

				if !more {
					break
				}
				frame, more = frames.Next()
			}
		}
	}

	// Write the newline character to the buffer.
	_ = enc.WriteByte('\n')
}
