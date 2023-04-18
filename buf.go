package olog

// Simple byte buffer for marshaling data.

import (
	"errors"
	"sync"
	"unicode/utf8"
)

const smallBufferSize = 64
const maxInt = int(^uint(0) >> 1)

// Declare a new sync.Pool object, which allows for efficient
// re-use of objects across goroutines.
var bufPool = sync.Pool{
	New: func() any {
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
	b.buf = utf8.AppendRune(b.buf[:m], r)
	return len(b.buf) - m, nil
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
