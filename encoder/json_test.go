package encoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

type testFormatter struct {
	s string
}

func (t *testFormatter) Format(f fmt.State, c rune) {
	if c == 'v' || c == 's' {
		_, _ = f.Write([]byte(t.s))
	}
}

type testStringer struct {
	s string
}

func (t *testStringer) String() string {
	return t.s
}

func TestJsonEncoder_WriteValue(t *testing.T) {
	tests := []struct {
		value interface{}
		want  string
	}{
		{
			value: "hello world",
			want:  "\"hello world\"",
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
			want:  "\"aGVsbG8gXCB3b3JsZA==\"",
		},
		{
			value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
			want:  "\"2020-01-01T00:00:00+08:00\"",
		},
		{
			value: errors.New("haha\""),
			want:  "\"haha\\\"\"",
		},
		{
			value: &testFormatter{s: "hello \\ world"},
			want:  "\"hello \\\\ world\"",
		},
		{
			value: &testStringer{s: "{\"name\": \"test\"}"},
			want:  "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			value: map[string]string{"name": "lisi", "age": "18"},
			want:  "\"{\\\"age\\\":\\\"18\\\",\\\"name\\\":\\\"lisi\\\"}\"",
		},
		{
			value: []int{1, 4, 10},
			want:  "\"[1,4,10]\"",
		},
	}

	e := JsonEncoder{&Buffer{}}

	for _, tt := range tests {
		e.Reset()
		e.WriteValue(tt.value)
		if string(e.Bytes()) != tt.want {
			t.Fatalf("get %s, want %s", string(e.Bytes()), tt.want)
		}
	}
}

func TestJsonEncoder_WriteEscapedString(t *testing.T) {
	tests := []struct {
		s  string
		es string
	}{
		{
			s:  "hello world",
			es: `hello world`,
		},
		{
			s:  "what happen\n",
			es: "what happen\\n",
		},
		{
			s:  "what \\ happen",
			es: "what \\\\ happen",
		},
		{
			s:  "{\"name\": \"test\"}",
			es: "{\\\"name\\\": \\\"test\\\"}",
		},
		{
			s:  "{\"name\": \"test\", \"age\": 18}",
			es: "{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}",
		},
	}

	e := JsonEncoder{&Buffer{}}
	for _, tt := range tests {
		e.Reset()
		e.WriteEscapedString(tt.s)
		if tt.es != string(e.Bytes()) {
			t.Fatalf("get %s, want %s", string(e.Bytes()), tt.es)
		}
	}

	str := `{"name": "test", "age": 18, "success": true}`
	e.Reset()
	e.WriteQuote()
	e.WriteEscapedString(str)
	e.WriteQuote()
	s, err := strconv.Unquote(string(e.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatal(err)
	}

	if m["name"] != "test" {
		t.Fatalf("get %s, want %s", m["name"], "test")
	}
}

func TestJsonEncoder_Write(t *testing.T) {
	tests := []struct {
		s  []byte
		es []byte
	}{
		{
			s:  []byte("hello world"),
			es: []byte(`hello world`),
		},
		{
			s:  []byte("what happen\n"),
			es: []byte("what happen\\n"),
		},
		{
			s:  []byte("what \\ happen"),
			es: []byte("what \\\\ happen"),
		},
		{
			s:  []byte("{\"name\": \"test\"}"),
			es: []byte("{\\\"name\\\": \\\"test\\\"}"),
		},
		{
			s:  []byte("{\"name\": \"test\", \"age\": 18}"),
			es: []byte("{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}"),
		},
		{
			s:  []byte{0, '\n', 1, 5},
			es: []byte("\\u0000\\n\\u0001\\u0005"),
		},
	}

	e := JsonEncoder{&Buffer{}}
	for _, tt := range tests {
		e.Reset()
		e.Write(tt.s)
		if string(tt.es) != string(e.Bytes()) {
			t.Fatalf("get %s, want %s", string(e.Bytes()), string(tt.es))
		}
	}

	bs := []byte(`{"name": "test", "age": 18, "success": true}`)
	e.Reset()
	e.WriteQuote()
	e.Write(bs)
	e.WriteQuote()
	s, err := strconv.Unquote(string(e.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatal(err)
	}

	if m["name"] != "test" {
		t.Fatalf("get %s, want %s", m["name"], "test")
	}
}
