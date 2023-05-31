package encoder

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
	"unicode/utf8"
)

const (
	hex   = "0123456789abcdef"
	quote = '"'
)

type JsonEncoder struct {
	*Buffer
}

func (j JsonEncoder) WriteNull() {
	_, _ = j.WriteString("null")
}

func (j JsonEncoder) WriteFloat(n float64, bitSize int) {
	switch {
	case math.IsNaN(n):
		_, _ = j.WriteString(`"NaN"`)
	case math.IsInf(n, +1):
		_, _ = j.WriteString(`"Infinity"`)
	case math.IsInf(n, -1):
		_, _ = j.WriteString(`"-Infinity"`)
	default:
		// JSON number formatting logic based on encoding/json.
		// See floatEncoder.encode for reference.
		f := byte('f')
		if abs := math.Abs(n); abs != 0 {
			if bitSize == 64 && (abs < 1e-6 || abs >= 1e21) ||
				bitSize == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
				f = 'e'
			}
		}
		j.Buffer.WriteFloat(n, f, bitSize)
	}
}

func (j JsonEncoder) StartObject() {
	_ = j.WriteByte('{')
}

func (j JsonEncoder) EndObject() {
	_ = j.WriteByte('}')
}

func (j JsonEncoder) StartArray() {
	_ = j.WriteByte('[')
}

func (j JsonEncoder) EndArray() {
	_ = j.WriteByte(']')
}

func (j JsonEncoder) WriteSeparator() {
	_ = j.WriteByte(',')
}

func (j JsonEncoder) WriteQuote() {
	_ = j.WriteByte(quote)
}

func (j JsonEncoder) WriteName(s string) {
	_ = j.WriteByte(quote)
	j.WriteEscapedString(s)
	_, _ = j.WriteString(`":`)
}

func (j JsonEncoder) WriteValue(value any) {
	switch v := value.(type) {
	case string:
		_ = j.WriteByte(quote)
		j.WriteEscapedString(v)
		_ = j.WriteByte(quote)
	case []byte:
		_ = j.WriteByte(quote)
		j.WriteBase64(v)
		_ = j.WriteByte(quote)
	case error:
		_ = j.WriteByte(quote)
		j.WriteEscapedString(v.Error())
		_ = j.WriteByte(quote)
	case time.Time:
		_ = j.WriteByte(quote)
		j.WriteTime(v, time.RFC3339)
		_ = j.WriteByte(quote)
	case nil:
		j.WriteNull()
	case int:
		j.WriteInt64(int64(v))
	case int8:
		j.WriteInt64(int64(v))
	case int16:
		j.WriteInt64(int64(v))
	case int32:
		j.WriteInt64(int64(v))
	case int64:
		j.WriteInt64(v)
	case uint:
		j.WriteUint64(uint64(v))
	case uint8:
		j.WriteUint64(uint64(v))
	case uint16:
		j.WriteUint64(uint64(v))
	case uint32:
		j.WriteUint64(uint64(v))
	case uint64:
		j.WriteUint64(v)
	case float32:
		j.WriteFloat(float64(v), 32)
	case float64:
		j.WriteFloat(v, 64)
	case bool:
		j.WriteBool(v)
	case fmt.Formatter:
		_ = j.WriteByte(quote)
		v.Format(j, 'v')
		_ = j.WriteByte(quote)
	case fmt.Stringer:
		_ = j.WriteByte(quote)
		j.WriteEscapedString(v.String())
		_ = j.WriteByte(quote)
	default:
		_ = j.WriteByte(quote)
		enc := json.NewEncoder(j)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(value); err != nil {
			j.WriteEscapedString("json.Marshal err: " + err.Error())
			_ = j.WriteByte(quote)
			return
		}

		// drop \n
		j.DropTail(2)
		_ = j.WriteByte(quote)
	}
}

func (j JsonEncoder) Write(s []byte) (n int, err error) {
	l := j.Len()
	j.growCap(len(s) + 8)

	start := 0
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if SafeSet[bt] {
				i++
				continue
			}
			if start < i {
				_, _ = j.Buffer.Write(s[start:i])
			}
			_ = j.WriteByte('\\')
			switch bt {
			case '\\', '"':
				_ = j.WriteByte(bt)
			case '\b':
				_ = j.WriteByte('b')
			case '\f':
				_ = j.WriteByte('f')
			case '\n':
				_ = j.WriteByte('n')
			case '\r':
				_ = j.WriteByte('r')
			case '\t':
				_ = j.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				_, _ = j.WriteString(`u00`)
				_ = j.WriteByte(hex[bt>>4])
				_ = j.WriteByte(hex[bt&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = j.Buffer.Write(s[start:i])
			}
			_, _ = j.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		i += size
	}

	if start < len(s) {
		_, _ = j.Buffer.Write(s[start:])
	}

	return j.Len() - l, nil
}

func (j JsonEncoder) Width() (int, bool) {
	return 0, false
}

func (j JsonEncoder) Precision() (int, bool) {
	return 0, false
}

func (j JsonEncoder) Flag(c int) bool {
	return false
}

func (j JsonEncoder) WriteEscapedString(s string) {
	j.growCap(len(s) + 8)

	start := 0
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if SafeSet[bt] {
				i++
				continue
			}
			if start < i {
				_, _ = j.WriteString(s[start:i])
			}
			_ = j.WriteByte('\\')
			switch bt {
			case '\\', '"':
				_ = j.WriteByte(bt)
			case '\b':
				_ = j.WriteByte('b')
			case '\f':
				_ = j.WriteByte('f')
			case '\n':
				_ = j.WriteByte('n')
			case '\r':
				_ = j.WriteByte('r')
			case '\t':
				_ = j.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				_, _ = j.WriteString(`u00`)
				_ = j.WriteByte(hex[bt>>4])
				_ = j.WriteByte(hex[bt&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = j.WriteString(s[start:i])
			}
			_, _ = j.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		_, _ = j.WriteString(s[start:])
	}
}

var SafeSet = [utf8.RuneSelf]bool{
	// 0 ~ 31 control characters, default to false
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
