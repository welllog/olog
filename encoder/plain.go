package encoder

import (
	"fmt"
	"math"
	"time"
)

type PlainEncoder struct {
	*Buffer
}

func (p PlainEncoder) WriteNull() {
	_, _ = p.WriteString("null")
}

func (p PlainEncoder) WriteFloat(n float64, bitSize int) {
	switch {
	case math.IsNaN(n):
		_, _ = p.WriteString(`"NaN"`)
	case math.IsInf(n, +1):
		_, _ = p.WriteString(`"Infinity"`)
	case math.IsInf(n, -1):
		_, _ = p.WriteString(`"-Infinity"`)
	default:
		p.Buffer.WriteFloat(n, 'f', bitSize)
	}
}

func (p PlainEncoder) WriteSeparator() {
	_ = p.WriteByte('\t')
}

func (p PlainEncoder) WriteName(s string) {
	_, _ = p.WriteString(s)
	_ = p.WriteByte('=')
}

func (p PlainEncoder) WriteValue(value any) {
	switch v := value.(type) {
	case string:
		_, _ = p.WriteString(v)
	case []byte:
		p.WriteBase64(v)
	case error:
		_, _ = p.WriteString(v.Error())
	case time.Time:
		p.WriteTime(v, time.RFC3339)
	case nil:
		p.WriteNull()
	case int:
		p.WriteInt64(int64(v))
	case int8:
		p.WriteInt64(int64(v))
	case int16:
		p.WriteInt64(int64(v))
	case int32:
		p.WriteInt64(int64(v))
	case int64:
		p.WriteInt64(v)
	case uint:
		p.WriteUint64(uint64(v))
	case uint8:
		p.WriteUint64(uint64(v))
	case uint16:
		p.WriteUint64(uint64(v))
	case uint32:
		p.WriteUint64(uint64(v))
	case uint64:
		p.WriteUint64(v)
	case float32:
		p.WriteFloat(float64(v), 32)
	case float64:
		p.WriteFloat(v, 64)
	case bool:
		p.WriteBool(v)
	case fmt.Formatter:
		v.Format(p, 'v')
	case fmt.Stringer:
		_, _ = p.WriteString(v.String())
	default:
		_, _ = fmt.Fprintf(p, "%+v", value)
	}
}

func (p PlainEncoder) Width() (int, bool) {
	return 0, false
}

func (p PlainEncoder) Precision() (int, bool) {
	return 0, false
}

func (p PlainEncoder) Flag(c int) bool {
	return false
}
