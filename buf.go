package olog

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	smallBufferSize = 64
	maxInt          = int(^uint(0) >> 1)
	lowerhex        = "0123456789abcdef"
	quote           = '"'
)

// Declare a new sync.Pool object, which allows for efficient
// re-use of objects across goroutines.
var bufPool = sync.Pool{
	New: func() interface{} {
		var b [256]byte
		return NewBuffer(b[:0])
	},
}

// getBuf to retrieve a *Buffer from the pool.
func getBuf() *Buffer {
	return bufPool.Get().(*Buffer)
}

// putBuf to return a *Buffer to the pool.
func putBuf(buf *Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

// Buffer byte buffer for marshaling data.
type Buffer struct {
	buf []byte
}

// NewBuffer creates and initializes a new Buffer using buf as its initial contents.
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }

func (b *Buffer) Len() int { return len(b.buf) }

func (b *Buffer) Cap() int { return cap(b.buf) }

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func (b *Buffer) Back(n int) {
	l := len(b.buf)
	if l < n {
		n = l
	}
	b.buf = b.buf[:l-n]
}

func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
}

func (b *Buffer) Swap(buf []byte) []byte {
	b.buf, buf = buf, b.buf
	return buf
}

func (b *Buffer) Grow(n int) {
	if n <= 0 {
		return
	}

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m]
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(s))
	if !ok {
		m = b.grow(len(s))
	}
	return copy(b.buf[m:], s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	m, ok := b.tryGrowByReslice(1)
	if !ok {
		m = b.grow(1)
	}
	b.buf[m] = c
	return nil
}

func (b *Buffer) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		_ = b.WriteByte(byte(r))
		return 1, nil
	}
	m, ok := b.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = b.grow(utf8.UTFMax)
	}
	n = utf8.EncodeRune(b.buf[m:m+utf8.UTFMax], r)
	b.buf = b.buf[:m+n]
	return n, nil
}

func (b *Buffer) WriteTime(t time.Time, layout string) {
	b.buf = t.AppendFormat(b.buf, layout)
}

func (b *Buffer) WriteInt64(n int64) {
	b.buf = strconv.AppendInt(b.buf, n, 10)
}

func (b *Buffer) WriteUint64(n uint64) {
	b.buf = strconv.AppendUint(b.buf, n, 10)
}

func (b *Buffer) WriteFloat64(f float64) {
	b.buf = strconv.AppendFloat(b.buf, f, 'f', -1, 64)
}

func (b *Buffer) WriteFloat32(f float32) {
	b.buf = strconv.AppendFloat(b.buf, float64(f), 'f', -1, 32)
}

func (b *Buffer) WriteBool(v bool) {
	b.buf = strconv.AppendBool(b.buf, v)
}

func (b *Buffer) WriteSprint(args ...interface{}) {
	_, _ = fmt.Fprint(b, args...)
}

func (b *Buffer) WriteSprintf(format string, args ...interface{}) {
	if len(args) == 0 {
		_, _ = b.WriteString(format)
		return
	}

	if format == "" {
		b.WriteSprint(args...)
		return
	}

	b.writeSprintf(format, args...)
}

func (b *Buffer) WriteQuoteSprint(args ...interface{}) {
	_ = b.WriteByte(quote)
	_, _ = fmt.Fprint(escapedWriter{buf: b}, args...)
	_ = b.WriteByte(quote)
}

func (b *Buffer) WriteQuoteSprintf(format string, args ...interface{}) {
	if len(args) == 0 {
		b.WriteQuoteString(format)
		return
	}

	if format == "" {
		b.WriteQuoteSprint(args...)
		return
	}

	b.writeQuoteSprintf(format, args...)
}

func (b *Buffer) WriteQuoteString(s string) {
	n := len(s) + 16
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m+1]
	b.buf[m] = quote

	b.WriteEscapedString(s)

	_ = b.WriteByte(quote)
}

func (b *Buffer) WriteQuoteBytes(s []byte) {
	n := len(s) + 16
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m+1]
	b.buf[m] = quote

	b.WriteEscapedBytes(s)

	_ = b.WriteByte(quote)
}

func (b *Buffer) WriteEscapedString(s string) {
	n := len(s) + 14
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m]

	start := 0
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if safeSet[bt] {
				i++
				continue
			}
			if start < i {
				_, _ = b.WriteString(s[start:i])
			}
			_ = b.WriteByte('\\')
			switch bt {
			case '\\', '"':
				_ = b.WriteByte(bt)
			case '\n':
				_ = b.WriteByte('n')
			case '\r':
				_ = b.WriteByte('r')
			case '\t':
				_ = b.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				_, _ = b.WriteString(`u00`)
				_ = b.WriteByte(lowerhex[bt>>4])
				_ = b.WriteByte(lowerhex[bt&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = b.WriteString(s[start:i])
			}
			_, _ = b.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				_, _ = b.WriteString(s[start:i])
			}
			_, _ = b.WriteString(`\u202`)
			_ = b.WriteByte(lowerhex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		_, _ = b.WriteString(s[start:])
	}
}

func (b *Buffer) WriteEscapedBytes(s []byte) {
	n := len(s) + 14
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m]

	start := 0
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if safeSet[bt] {
				i++
				continue
			}
			if start < i {
				_, _ = b.Write(s[start:i])
			}
			_ = b.WriteByte('\\')
			switch bt {
			case '\\', '"':
				_ = b.WriteByte(bt)
			case '\n':
				_ = b.WriteByte('n')
			case '\r':
				_ = b.WriteByte('r')
			case '\t':
				_ = b.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				_, _ = b.WriteString(`u00`)
				_ = b.WriteByte(lowerhex[bt>>4])
				_ = b.WriteByte(lowerhex[bt&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = b.Write(s[start:i])
			}
			_, _ = b.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				_, _ = b.Write(s[start:i])
			}
			_, _ = b.WriteString(`\u202`)
			_ = b.WriteByte(lowerhex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		_, _ = b.Write(s[start:])
	}
}

func (b *Buffer) WriteEscapedRune(r rune) {
	if uint32(r) < utf8.RuneSelf {
		if safeSet[byte(r)] {
			_ = b.WriteByte(byte(r))
			return
		}

		_ = b.WriteByte('\\')
		switch r {
		case '\\', '"':
			_ = b.WriteByte(byte(r))
		case '\n':
			_ = b.WriteByte('n')
		case '\r':
			_ = b.WriteByte('r')
		case '\t':
			_ = b.WriteByte('t')
		default:
			// This encodes bytes < 0x20 except for \t, \n and \r.
			// If escapeHTML is set, it also escapes <, >, and &
			// because they can lead to security holes when
			// user-controlled strings are rendered into JSON
			// and served to some browsers.
			_, _ = b.WriteString(`u00`)
			_ = b.WriteByte(lowerhex[r>>4])
			_ = b.WriteByte(lowerhex[r&0xF])
		}
		return
	}

	if r == utf8.RuneError {
		_, _ = b.WriteString(`\ufffd`)
		return
	}

	if r == '\u2028' || r == '\u2029' {
		_, _ = b.WriteString(`\u202`)
		_ = b.WriteByte(lowerhex[r&0xF])
		return
	}

	m, ok := b.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = b.grow(utf8.UTFMax)
	}
	n := utf8.EncodeRune(b.buf[m:m+utf8.UTFMax], r)
	b.buf = b.buf[:m+n]
}

func (b *Buffer) Width() (int, bool) {
	return 0, false
}

func (b *Buffer) Precision() (int, bool) {
	return 0, false
}

func (b *Buffer) Flag(c int) bool {
	return false
}

func (b *Buffer) writeSprintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(b, format, args...)
}

func (b *Buffer) writeQuoteSprintf(format string, args ...interface{}) {
	_ = b.WriteByte(quote)
	_, _ = fmt.Fprintf(escapedWriter{buf: b}, format, args...)
	_ = b.WriteByte(quote)
}

// tryGrowByReslice is a inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) grow(n int) int {
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}

	m := len(b.buf)
	c := cap(b.buf)
	if c > maxInt-c-n {
		panic(ErrTooLarge)
	}

	b.buf = growSlice(b.buf, n)
	b.buf = b.buf[:m+n]
	return m
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}

type escapedWriter struct {
	buf *Buffer
}

func (e escapedWriter) Write(p []byte) (n int, err error) {
	l := len(e.buf.buf)
	e.buf.WriteEscapedBytes(p)
	return len(e.buf.buf) - l, nil
}

func (e escapedWriter) Width() (int, bool) {
	return 0, false
}

func (e escapedWriter) Precision() (int, bool) {
	return 0, false
}

func (e escapedWriter) Flag(c int) bool {
	return false
}

var safeSet = [utf8.RuneSelf]bool{
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
