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
	buf.ResetBuffer(time.Now().AppendFormat(buf.Bytes(), o.timeFormat))
	_, _ = buf.WriteString(`","level":"`)
	_, _ = buf.WriteString(o.tag)
	if o.enableCaller {
		_, _ = buf.WriteString(`","caller":"`)
		_, _ = buf.WriteString(o.file)
		_ = buf.WriteByte(':')
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(o.line), 10))
	}
	_, _ = buf.WriteString(`","content":`)

	l := len(buf.buf)
	switch o.msgType {
	case msgTypePrint:
		buf.ResetBuffer(fmt.Append(buf.Bytes(), o.msgArgs...))
	case msgTypePrintf:
		buf.ResetBuffer(fmt.Appendf(buf.Bytes(), o.msgOrFormat, o.msgArgs...))
	case msgTypePrintMsg:
		_, _ = buf.WriteString(o.msgOrFormat)
	}

	quoteStr := string(buf.buf[l:])
	buf.buf = buf.buf[:l]
	buf.ResetBuffer(strconv.AppendQuote(buf.Bytes(), quoteStr))

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
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

// plainEncode to encode a logOption object as plain text and write it to the given Writer.
// Takes a pointer to the logOption object and a Writer object as input.
func plainEncode(o *logOption, w Writer) {
	buf := getBuf()
	buf.ResetBuffer(time.Now().AppendFormat(buf.Bytes(), o.timeFormat))
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
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(o.line), 10))
		_ = buf.WriteByte(sep)
	}

	switch o.msgType {
	case msgTypePrint:
		buf.ResetBuffer(fmt.Append(buf.Bytes(), o.msgArgs...))
	case msgTypePrintf:
		buf.ResetBuffer(fmt.Appendf(buf.Bytes(), o.msgOrFormat, o.msgArgs...))
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
	_, _ = w.Write(o.level, buf.Bytes())

	// Reset the buffer and return it to the pool for re-use.
	putBuf(buf)
}

func writeJsonValue(value any, buf *Buffer, quote bool) {
	switch v := value.(type) {
	case string:
		if quote {
			buf.ResetBuffer(strconv.AppendQuote(buf.Bytes(), v))
			return
		}
		_, _ = buf.WriteString(v)
	case fmt.Stringer:
		if quote {
			buf.ResetBuffer(strconv.AppendQuote(buf.Bytes(), v.String()))
			return
		}
		_, _ = buf.WriteString(v.String())
	case []byte:
		if quote {
			buf.ResetBuffer(strconv.AppendQuote(buf.Bytes(), *(*string)(unsafe.Pointer(&v))))
			return
		}
		_, _ = buf.Write(v)
	case nil:
		_, _ = buf.WriteString("null")
	case int:
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(v), 10))
	case int8:
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(v), 10))
	case int16:
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(v), 10))
	case int32:
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), int64(v), 10))
	case int64:
		buf.ResetBuffer(strconv.AppendInt(buf.Bytes(), v, 10))
	case uint:
		buf.ResetBuffer(strconv.AppendUint(buf.Bytes(), uint64(v), 10))
	case uint8:
		buf.ResetBuffer(strconv.AppendUint(buf.Bytes(), uint64(v), 10))
	case uint16:
		buf.ResetBuffer(strconv.AppendUint(buf.Bytes(), uint64(v), 10))
	case uint32:
		buf.ResetBuffer(strconv.AppendUint(buf.Bytes(), uint64(v), 10))
	case uint64:
		buf.ResetBuffer(strconv.AppendUint(buf.Bytes(), v, 10))
	case float32:
		buf.ResetBuffer(strconv.AppendFloat(buf.Bytes(), float64(v), 'f', -1, 32))
	case float64:
		buf.ResetBuffer(strconv.AppendFloat(buf.Bytes(), v, 'f', -1, 64))
	case bool:
		buf.ResetBuffer(strconv.AppendBool(buf.Bytes(), v))
	default:
		if quote {
			l := len(buf.buf)
			buf.ResetBuffer(fmt.Appendf(buf.Bytes(), "%+v", value))
			quoteStr := string(buf.buf[l:])
			buf.buf = buf.buf[:l]
			buf.ResetBuffer(strconv.AppendQuote(buf.Bytes(), quoteStr))
			return
		}
		buf.ResetBuffer(fmt.Appendf(buf.Bytes(), "%+v", value))
	}
}
