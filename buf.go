package olog

// Simple byte buffer for marshaling data.

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"
)

const smallBufferSize = 64
const maxInt = int(^uint(0) >> 1)
const lowerhex = "0123456789abcdef"

// Declare a new sync.Pool object, which allows for efficient
// re-use of objects across goroutines.
var bufPool = sync.Pool{
	New: func() interface{} {
		var b [200]byte
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

type Buffer struct {
	buf []byte
}

// NewBuffer creates and initializes a new Buffer using buf as its initial contents.
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }

func (b *Buffer) Len() int { return len(b.buf) }

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

func (b *Buffer) WriteSprint(args ...interface{}) {
	_, _ = fmt.Fprint(b, args...)
}

func (b *Buffer) WriteSprintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(b, format, args...)
}

func (b *Buffer) WriteQuoteSprint(args ...interface{}) {
	l := len(b.buf)
	_, _ = fmt.Fprint(b, args...)
	s := string(b.buf[l:])
	b.buf = b.buf[:l]
	b.WriteQuoteString(s)
}

func (b *Buffer) WriteQuoteSprintf(format string, args ...interface{}) {
	l := len(b.buf)
	_, _ = fmt.Fprintf(b, format, args...)
	s := string(b.buf[l:])
	b.buf = b.buf[:l]
	b.WriteQuoteString(s)
}

func (b *Buffer) WriteAny(value interface{}, quoteStr bool) {
	switch v := value.(type) {
	case string:
		if quoteStr {
			b.WriteQuoteString(v)
			return
		}
		_, _ = b.WriteString(v)
	case []byte:
		if quoteStr {
			b.WriteQuoteString(*(*string)(unsafe.Pointer(&v)))
			return
		}
		_, _ = b.Write(v)
	case error:
		if quoteStr {
			b.WriteQuoteString(v.Error())
			return
		}
		_, _ = b.WriteString(v.Error())
	case time.Time:
		if quoteStr {
			_ = b.WriteByte('"')
			b.buf = v.AppendFormat(b.buf, time.RFC3339)
			_ = b.WriteByte('"')
			return
		}
		b.buf = v.AppendFormat(b.buf, time.RFC3339)
	case nil:
		_, _ = b.WriteString("null")
	case int:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int8:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int16:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int32:
		b.buf = strconv.AppendInt(b.buf, int64(v), 10)
	case int64:
		b.buf = strconv.AppendInt(b.buf, v, 10)
	case uint:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint8:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint16:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint32:
		b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
	case uint64:
		b.buf = strconv.AppendUint(b.buf, v, 10)
	case float32:
		b.buf = strconv.AppendFloat(b.buf, float64(v), 'f', -1, 32)
	case float64:
		b.buf = strconv.AppendFloat(b.buf, v, 'f', -1, 64)
	case bool:
		b.buf = strconv.AppendBool(b.buf, v)
	case fmt.Stringer:
		if quoteStr {
			b.WriteQuoteString(v.String())
			return
		}
		_, _ = b.WriteString(v.String())
	default:
		if quoteStr {
			b.WriteQuoteSprintf("%+v", value)
			return
		}
		b.WriteSprintf("%+v", value)
	}
}

func (b *Buffer) WriteQuoteString(s string) {
	if cap(b.buf)-len(b.buf) < len(s) {
		b.buf = growSlice(b.buf, len(s)+2)
	}
	_ = b.WriteByte('"')
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			_, _ = b.WriteString(`\x`)
			_ = b.WriteByte(lowerhex[s[0]>>4])
			_ = b.WriteByte(lowerhex[s[0]&0xF])
			continue
		}
		b.WriteEscapedRune(r)
	}
	_ = b.WriteByte('"')
}

func (b *Buffer) WriteEscapedRune(r rune) {
	if r == '"' || r == '\\' { // always backslashed
		_ = b.WriteByte('\\')
		_ = b.WriteByte(byte(r))
		return
	}

	if strconv.IsPrint(r) {
		_, _ = b.WriteRune(r)
		return
	}

	switch r {
	case '\a':
		_, _ = b.WriteString(`\a`)
	case '\b':
		_, _ = b.WriteString(`\b`)
	case '\f':
		_, _ = b.WriteString(`\f`)
	case '\n':
		_, _ = b.WriteString(`\n`)
	case '\r':
		_, _ = b.WriteString(`\r`)
	case '\t':
		_, _ = b.WriteString(`\t`)
	case '\v':
		_, _ = b.WriteString(`\v`)
	default:
		switch {
		case r < ' ' || r == 0x7f:
			_, _ = b.WriteString(`\x`)
			_ = b.WriteByte(lowerhex[byte(r)>>4])
			_ = b.WriteByte(lowerhex[byte(r)&0xF])
		case !utf8.ValidRune(r):
			r = 0xFFFD
			fallthrough
		case r < 0x10000:
			_, _ = b.WriteString(`\u`)
			for s := 12; s >= 0; s -= 4 {
				_ = b.WriteByte(lowerhex[r>>uint(s)&0xF])
			}
		default:
			_, _ = b.WriteString(`\U`)
			for s := 28; s >= 0; s -= 4 {
				_ = b.WriteByte(lowerhex[r>>uint(s)&0xF])
			}
		}
	}
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
