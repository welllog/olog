package olog

import (
	"encoding/json"
	"fmt"
	"time"
)

type ValueEncoder struct {
	*Buffer
}

func (e *ValueEncoder) EncodeJsonValue(value interface{}) {
	switch v := value.(type) {
	case string:
		e.WriteQuoteString(v)
	case []byte:
		e.WriteQuoteBytes(v)
	case error:
		e.WriteQuoteString(v.Error())
	case time.Time:
		_ = e.WriteByte(quote)
		e.WriteTime(v, time.RFC3339)
		_ = e.WriteByte(quote)
	case nil:
		_, _ = e.WriteString("null")
	case int:
		e.WriteInt64(int64(v))
	case int8:
		e.WriteInt64(int64(v))
	case int16:
		e.WriteInt64(int64(v))
	case int32:
		e.WriteInt64(int64(v))
	case int64:
		e.WriteInt64(v)
	case uint:
		e.WriteUint64(uint64(v))
	case uint8:
		e.WriteUint64(uint64(v))
	case uint16:
		e.WriteUint64(uint64(v))
	case uint32:
		e.WriteUint64(uint64(v))
	case uint64:
		e.WriteUint64(v)
	case float32:
		e.WriteFloat32(v)
	case float64:
		e.WriteFloat64(v)
	case bool:
		e.WriteBool(v)
	case fmt.Formatter:
		_ = e.WriteByte(quote)
		v.Format(escapedWriter{buf: e.Buffer}, 'v')
		_ = e.WriteByte(quote)
	case fmt.Stringer:
		e.WriteQuoteString(v.String())
	default:
		_ = e.WriteByte(quote)
		enc := json.NewEncoder(escapedWriter{buf: e.Buffer})
		enc.SetEscapeHTML(false)

		if err := enc.Encode(value); err != nil {
			e.WriteEscapedString("json.Marshal err: " + err.Error())
			_ = e.WriteByte(quote)
			return
		}

		// drop \n
		e.Back(2)
		_ = e.WriteByte(quote)
	}
}

func (e *ValueEncoder) EncodeValue(value interface{}) {
	switch v := value.(type) {
	case string:
		_, _ = e.WriteString(v)
	case []byte:
		_, _ = e.Write(v)
	case error:
		_, _ = e.WriteString(v.Error())
	case time.Time:
		e.WriteTime(v, time.RFC3339)
	case nil:
		_, _ = e.WriteString("null")
	case int:
		e.WriteInt64(int64(v))
	case int8:
		e.WriteInt64(int64(v))
	case int16:
		e.WriteInt64(int64(v))
	case int32:
		e.WriteInt64(int64(v))
	case int64:
		e.WriteInt64(v)
	case uint:
		e.WriteUint64(uint64(v))
	case uint8:
		e.WriteUint64(uint64(v))
	case uint16:
		e.WriteUint64(uint64(v))
	case uint32:
		e.WriteUint64(uint64(v))
	case uint64:
		e.WriteUint64(v)
	case float32:
		e.WriteFloat32(v)
	case float64:
		e.WriteFloat64(v)
	case bool:
		e.WriteBool(v)
	case fmt.Formatter:
		v.Format(e.Buffer, 'v')
	case fmt.Stringer:
		_, _ = e.WriteString(v.String())
	default:
		_, _ = fmt.Fprintf(e.Buffer, "%+v", value)
	}
}
