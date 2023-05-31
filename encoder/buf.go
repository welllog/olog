package encoder

import (
	"encoding/base64"
	"errors"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	smallBufferSize = 64
	maxInt          = int(^uint(0) >> 1)
)

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

func (b *Buffer) DropTail(n int) {
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

func (b *Buffer) WriteFloat(f float64, fmt byte, bitSize int) {
	l := len(b.buf)
	b.buf = strconv.AppendFloat(b.buf, f, fmt, -1, bitSize)
	if fmt == 'e' {
		n := len(b.buf)
		if n-l >= 4 && b.buf[n-4] == 'e' && b.buf[n-3] == '-' && b.buf[n-2] == '0' {
			b.buf[n-2] = b.buf[n-1]
			b.buf = b.buf[:n-1]
		}
	}
}

func (b *Buffer) WriteBool(v bool) {
	b.buf = strconv.AppendBool(b.buf, v)
}

func (b *Buffer) WriteBase64(p []byte) {
	n := base64.StdEncoding.EncodedLen(len(p))
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	base64.StdEncoding.Encode(b.buf[m:], p)
}

func (b *Buffer) growCap(n int) {
	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}
	b.buf = b.buf[:m]
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

	c := cap(b.buf)
	if c > maxInt-c-n {
		panic(ErrTooLarge)
	}

	l := len(b.buf)
	nc := l + n
	if nc < 2*c {
		nc = 2 * c
	} else {
		nc += smallBufferSize
	}

	buf := make([]byte, nc)
	copy(buf, b.buf)
	b.buf = buf[:l+n]

	return l
}
