package olog

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
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
	buf.buf = time.Now().AppendFormat(buf.buf, o.timeFormat)
	_, _ = buf.WriteString(`","level":"`)
	_, _ = buf.WriteString(o.tag)
	if o.enableCaller {
		_, _ = buf.WriteString(`","caller":"`)
		_, _ = buf.WriteString(o.file)
		_ = buf.WriteByte(':')
		buf.buf = strconv.AppendInt(buf.buf, int64(o.line), 10)
	}
	_, _ = buf.WriteString(`","content":`)

	l := len(buf.buf)
	switch o.msgType {
	case msgTypePrint:
		buf.buf = fmt.Append(buf.buf, o.msgArgs...)
	case msgTypePrintf:
		buf.buf = fmt.Appendf(buf.buf, o.msgOrFormat, o.msgArgs...)
	case msgTypePrintMsg:
		_, _ = buf.WriteString(o.msgOrFormat)
	}

	quoteStr := string(buf.buf[l:])
	buf.buf = buf.buf[:l]
	buf.buf = strconv.AppendQuote(buf.buf, quoteStr)

	// Loop over the fields of the logOption object and write them to the buffer as JSON.
	for _, field := range o.fields {
		// Write the field key and value to the buffer as JSON.
		_, _ = buf.WriteString(`,"`)
		_, _ = buf.WriteString(field.Key)
		_, _ = buf.WriteString(`":`)
		writeJsonValue(field.Value, buf, true)
	}

	// Write the closing curly brace and newline character to the buffer.
	_, _ = buf.WriteString("}\n")

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.buf)

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

// plainEncode to encode a logOption object as plain text and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func plainEncode(o *logOption, w Writer) {
	buf := getBuf()
	buf.buf = time.Now().AppendFormat(buf.buf, o.timeFormat)
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
		buf.buf = strconv.AppendInt(buf.buf, int64(o.line), 10)
		_ = buf.WriteByte(sep)
	}

	switch o.msgType {
	case msgTypePrint:
		buf.buf = fmt.Append(buf.buf, o.msgArgs...)
	case msgTypePrintf:
		buf.buf = fmt.Appendf(buf.buf, o.msgOrFormat, o.msgArgs...)
	case msgTypePrintMsg:
		_, _ = buf.WriteString(o.msgOrFormat)
	}

	// Loop over the fields of the logOption object and write them to the buffer as plain text.
	for _, field := range o.fields {
		_ = buf.WriteByte(sep)
		_, _ = buf.WriteString(field.Key)
		_ = buf.WriteByte(sep)
		writeJsonValue(field.Value, buf, false)
	}

	// Write the newline character to the buffer.
	_ = buf.WriteByte('\n')

	// Write the contents of the buffer to the Writer.
	_, _ = w.Write(o.level, buf.buf)

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

func writeJsonValue(value any, buf *Buffer, quote bool) {
	switch v := value.(type) {
	case string:
		if quote {
			buf.buf = strconv.AppendQuote(buf.buf, v)
			return
		}
		_, _ = buf.WriteString(v)
	case fmt.Stringer:
		if quote {
			buf.buf = strconv.AppendQuote(buf.buf, v.String())
			return
		}
		_, _ = buf.WriteString(v.String())
	case []byte:
		if quote {
			buf.buf = strconv.AppendQuote(buf.buf, *(*string)(unsafe.Pointer(&v)))
			return
		}
		_, _ = buf.Write(v)
	case error:
		if quote {
			buf.buf = strconv.AppendQuote(buf.buf, v.Error())
			return
		}
		buf.buf = append(buf.buf, v.Error()...)
	case time.Time:
		_ = buf.WriteByte('"')
		buf.buf = v.AppendFormat(buf.buf, time.RFC3339)
		_ = buf.WriteByte('"')
	case nil:
		_, _ = buf.WriteString("null")
	case int:
		buf.buf = strconv.AppendInt(buf.buf, int64(v), 10)
	case int8:
		buf.buf = strconv.AppendInt(buf.buf, int64(v), 10)
	case int16:
		buf.buf = strconv.AppendInt(buf.buf, int64(v), 10)
	case int32:
		buf.buf = strconv.AppendInt(buf.buf, int64(v), 10)
	case int64:
		buf.buf = strconv.AppendInt(buf.buf, v, 10)
	case uint:
		buf.buf = strconv.AppendUint(buf.buf, uint64(v), 10)
	case uint8:
		buf.buf = strconv.AppendUint(buf.buf, uint64(v), 10)
	case uint16:
		buf.buf = strconv.AppendUint(buf.buf, uint64(v), 10)
	case uint32:
		buf.buf = strconv.AppendUint(buf.buf, uint64(v), 10)
	case uint64:
		buf.buf = strconv.AppendUint(buf.buf, v, 10)
	case float32:
		buf.buf = strconv.AppendFloat(buf.buf, float64(v), 'f', -1, 32)
	case float64:
		buf.buf = strconv.AppendFloat(buf.buf, v, 'f', -1, 64)
	case bool:
		buf.buf = strconv.AppendBool(buf.buf, v)
	default:
		if quote {
			l := len(buf.buf)
			buf.buf = fmt.Appendf(buf.buf, "%+v", value)
			quoteStr := string(buf.buf[l:])
			buf.buf = buf.buf[:l]
			buf.buf = strconv.AppendQuote(buf.buf, quoteStr)
			return
		}
		buf.buf = fmt.Appendf(buf.buf, "%+v", value)
	}
}
