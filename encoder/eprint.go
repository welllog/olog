package encoder

import (
	"fmt"
	"io"
	"unsafe"
)

func EPrint(w io.Writer, args ...any) (n int, err error) {
	if len(args) == 0 {
		return 0, nil
	}
	if len(args) == 1 {
		if s, ok := args[0].(string); ok {
			return w.Write(*(*[]byte)(unsafe.Pointer(
				&struct {
					string
					Cap int
				}{s, len(s)},
			)))
		}
	}
	return fmt.Fprint(w, args...)
}

func EPrintf(w io.Writer, format string, args ...any) (n int, err error) {
	if len(args) == 0 {
		if format == "" {
			return 0, nil
		}

		return w.Write(*(*[]byte)(unsafe.Pointer(
			&struct {
				string
				Cap int
			}{format, len(format)},
		)))
	}

	if format == "" {
		if len(args) == 1 {
			if s, ok := args[0].(string); ok {
				return w.Write(*(*[]byte)(unsafe.Pointer(
					&struct {
						string
						Cap int
					}{s, len(s)},
				)))
			}
		}
		return fmt.Fprint(w, args...)
	}

	return fmt.Fprintf(w, format, args...)
}
