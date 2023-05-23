package olog

import (
	"encoding/json"
	"fmt"
	"time"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

type encoder struct {
	*Buffer
}

func (e *encoder) EncodeJsonValue(value interface{}) {
	switch v := value.(type) {
	case string:
		e.WriteQuoteString(v)
	case []byte:
		e.WriteQuoteBytes(v)
	case error:
		e.WriteQuoteString(v.Error())
	case time.Time:
		_ = e.WriteByte('"')
		e.WriteTime(v, time.RFC3339)
		_ = e.WriteByte('"')
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
	case fmt.Stringer:
		e.WriteQuoteString(v.String())
	case Marshaler:
		b, err := v.Marshal()
		if err != nil {
			e.WriteQuoteString("olog.Marshaler marshal err: " + err.Error())
		} else {
			e.WriteQuoteBytes(b)
		}
	default:
		b, err := json.Marshal(value)
		if err != nil {
			e.WriteQuoteString("json.Marshal err: " + err.Error())
		} else {
			e.WriteQuoteBytes(b)
		}
	}
}

func (e *encoder) EncodeValue(v interface{}) {

}
