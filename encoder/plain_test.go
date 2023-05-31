package encoder

import (
	"errors"
	"testing"
	"time"
)

func TestPlainEncoder_WriteValue(t *testing.T) {
	tests := []struct {
		value any
		want  string
	}{
		{
			value: "hello world",
			want:  "hello world",
		},
		{
			value: 123,
			want:  "123",
		},
		{
			value: 11.2,
			want:  "11.2",
		},
		{
			value: true,
			want:  "true",
		},
		{
			value: false,
			want:  "false",
		},
		{
			value: uint8(1),
			want:  "1",
		},
		{
			value: []byte("hello \\ world"),
			want:  "aGVsbG8gXCB3b3JsZA==",
		},
		{
			value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			want:  "2020-01-01T00:00:00+08:00",
		},
		{
			value: errors.New("haha\""),
			want:  "haha\"",
		},
		{
			value: &testFormatter{s: "hello \\ world"},
			want:  "hello \\ world",
		},
		{
			value: &testStringer{s: "{\"name\": \"test\"}"},
			want:  "{\"name\": \"test\"}",
		},
		{
			value: map[string]string{"name": "lisi", "age": "18"},
			want:  "map[age:18 name:lisi]",
		},
	}
	e := PlainEncoder{&Buffer{}}

	for _, tt := range tests {
		e.Reset()
		e.WriteValue(tt.value)
		if string(e.Bytes()) != tt.want {
			t.Fatalf("get %s, want %s", string(e.Bytes()), tt.want)
		}
	}
}
