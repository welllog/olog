package olog

import (
	"runtime"
	"unicode/utf8"
	"unsafe"

	"github.com/welllog/olog/encoder"
)

const lowerhex = "0123456789abcdef"

// TrimLineEnding function removes the trailing newline character ('\n') from the end of the byte slice.
func TrimLineEnding(b []byte) []byte {
	if l := len(b); l > 0 && b[l-1] == '\n' {
		return b[:l-1]
	}
	return b
}

func getCaller(skip int8) (string, int) {
	_, file, line, ok := runtime.Caller(int(skip))
	if !ok {
		return "", 0
	}
	return file, line
}

func getCallerFrames(skip int8, size uint8) *runtime.Frames {
	pc := make([]uintptr, size)
	n := runtime.Callers(int(skip+1), pc)

	return runtime.CallersFrames(pc[:n])
}

func shortFile(file string) string {
	if file == "" {
		return "???"
	}

	var count int
	idx := -1
	for i := len(file) - 5; i >= 0; i-- {
		if file[i] == '/' {
			count++
			if count == 2 {
				idx = i
				break
			}
		}
	}
	if idx == -1 {
		return file
	}
	return file[idx+1:]
}

// EscapedString function returns a string with all the special characters escaped.
func EscapedString(s string) string {
	index := indexNeedEscapedString(s)
	if index == -1 {
		return s
	}

	buf := make([]byte, 0, len(s)+8)
	if index > 0 {
		buf = append(buf, s[:index]...)
	}
	buf = appendEscapedString(buf, s[index:])
	return *(*string)(unsafe.Pointer(&buf))
}

func appendEscapedString(dst []byte, s string) []byte {
	start := 0
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if encoder.SafeSet[bt] {
				i++
				continue
			}
			if start < i {
				dst = append(dst, s[start:i]...)
			}
			dst = append(dst, '\\')
			switch bt {
			case '\\', '"':
				dst = append(dst, bt)
			case '\b':
				dst = append(dst, 'b')
			case '\f':
				dst = append(dst, 'f')
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				dst = append(dst, `u00`...)
				dst = append(dst, lowerhex[bt>>4])
				dst = append(dst, lowerhex[bt&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				dst = append(dst, s[start:i]...)
			}
			dst = append(dst, `\ufffd`...)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}

	return dst
}

func indexNeedEscapedString(s string) int {
	for i := 0; i < len(s); {
		if bt := s[i]; bt < utf8.RuneSelf {
			if encoder.SafeSet[bt] {
				i++
				continue
			}
			return i
		}

		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			return i
		}
		i += size
	}
	return -1
}
