package olog

import (
	"time"
)

// EncodeType is an enumeration type for different encoding types.
type EncodeType int

const (
	// JSON represents the JSON encoding type.
	JSON EncodeType = iota
	// PLAIN represents the plain text encoding type.
	PLAIN
	// sep is the separator between fields in the log message.
	sep = '\t'
)

// jsonEncode to encode a logOption object as JSON and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func jsonEncode(o *logOption, w Writer) {
	buf := getBuf()

	_, _ = buf.WriteString(`{"@timestamp":"`)
	buf.WriteTime(time.Now(), o.timeFormat)
	_, _ = buf.WriteString(`","level":"`)
	_, _ = buf.WriteString(o.tag)
	if o.enableCaller {
		_, _ = buf.WriteString(`","caller":"`)
		_, _ = buf.WriteString(o.file)
		_ = buf.WriteByte(':')
		buf.WriteInt64(int64(o.line))
	}
	_, _ = buf.WriteString(`","content":`)

	switch o.msgType {
	case msgTypePrint:
		buf.WriteQuoteSprint(o.msgArgs...)
	case msgTypePrintf:
		buf.WriteQuoteSprintf(o.msgOrFormat, o.msgArgs...)
	case msgTypePrintMsg:
		buf.WriteQuoteString(o.msgOrFormat)
	}

	// Loop over the fields of the logOption object and write them to the buffer as JSON.
	for _, field := range o.fields {
		// Write the field key and value to the buffer as JSON.
		_, _ = buf.WriteString(`,"`)
		_, _ = buf.WriteString(field.Key)
		_, _ = buf.WriteString(`":`)
		buf.WriteAny(field.Value, true)
	}

	// Write the closing curly brace and newline character to the buffer.
	_, _ = buf.WriteString("}\n")

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

// plainEncode to encode a logOption object as plain text and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func plainEncode(o *logOption, w Writer) {
	buf := getBuf()

	buf.WriteTime(time.Now(), o.timeFormat)
	_ = buf.WriteByte(sep)
	if o.enableColor {
		writeLevelWithColor(o.level, o.tag, buf)
	} else {
		_, _ = buf.WriteString(o.tag)
	}
	_ = buf.WriteByte(sep)
	if o.enableCaller {
		_, _ = buf.WriteString(o.file)
		_ = buf.WriteByte(':')
		buf.WriteInt64(int64(o.line))
		_ = buf.WriteByte(sep)
	}

	switch o.msgType {
	case msgTypePrint:
		buf.WriteSprint(o.msgArgs...)
	case msgTypePrintf:
		buf.WriteSprintf(o.msgOrFormat, o.msgArgs...)
	case msgTypePrintMsg:
		_, _ = buf.WriteString(o.msgOrFormat)
	}

	// Loop over the fields of the logOption object and write them to the buffer as plain text.
	for _, field := range o.fields {
		_ = buf.WriteByte(sep)
		_, _ = buf.WriteString(field.Key)
		_ = buf.WriteByte('=')
		buf.WriteAny(field.Value, false)
	}

	// Write the newline character to the buffer.
	_ = buf.WriteByte('\n')

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}
