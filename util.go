package olog

import (
	"runtime"
)

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
