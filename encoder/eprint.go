package encoder

import (
	"fmt"
	"io"
	"unsafe"
)

func EPrint(w io.Writer, args ...interface{}) (n int, err error) {
	return fmt.Fprint(w, args...)
}

func EPrintf(w io.Writer, format string, args ...interface{}) (n int, err error) {
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
		return fmt.Fprint(w, args...)
	}

	return fmt.Fprintf(w, format, args...)
}
