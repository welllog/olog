package olog

import (
	"encoding/json"
	"strconv"
	"testing"
)

func TestBuffer_WriteQuoteString(t *testing.T) {
	tests := []struct {
		s  string
		es string
	}{
		{
			s:  "hello world",
			es: `"hello world"`,
		},
		{
			s:  "what happen\n",
			es: "\"what happen\\n\"",
		},
		{
			s:  "what \\ happen",
			es: "\"what \\\\ happen\"",
		},
		{
			s:  "{\"name\": \"test\"}",
			es: "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			s:  "{\"name\": \"test\", \"age\": 18}",
			es: "\"{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}\"",
		},
	}

	b := NewBuffer(nil)
	for _, tt := range tests {
		b.Reset()
		b.WriteQuoteString(tt.s)
		if tt.es != string(b.Bytes()) {
			t.Fatalf("get %s, want %s", string(b.Bytes()), tt.es)
		}
	}

	str := `{"name": "test", "age": 18, "success": true}`
	b.Reset()
	b.WriteQuoteString(str)
	s, err := strconv.Unquote(string(b.Bytes()))
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

func TestBuffer_WriteQuoteBytes(t *testing.T) {
	tests := []struct {
		s  []byte
		es []byte
	}{
		{
			s:  []byte("hello world"),
			es: []byte(`"hello world"`),
		},
		{
			s:  []byte("what happen\n"),
			es: []byte("\"what happen\\n\""),
		},
		{
			s:  []byte("what \\ happen"),
			es: []byte("\"what \\\\ happen\""),
		},
		{
			s:  []byte("{\"name\": \"test\"}"),
			es: []byte("\"{\\\"name\\\": \\\"test\\\"}\""),
		},
		{
			s:  []byte("{\"name\": \"test\", \"age\": 18}"),
			es: []byte("\"{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}\""),
		},
		{
			s:  []byte{0, '\n', 1, 5},
			es: []byte("\"\\u0000\\n\\u0001\\u0005\""),
		},
	}
	b := NewBuffer(nil)
	for _, tt := range tests {
		b.Reset()
		b.WriteQuoteBytes(tt.s)
		if string(tt.es) != string(b.Bytes()) {
			t.Fatalf("get %s, want %s", string(b.Bytes()), string(tt.es))
		}
	}

	bs := []byte(`{"name": "test", "age": 18, "success": true}`)
	b.Reset()
	b.WriteQuoteBytes(bs)
	s, err := strconv.Unquote(string(b.Bytes()))
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

func TestBuffer_WriteQuoteSprint(t *testing.T) {
	tests := []struct {
		a  []interface{}
		es string
	}{
		{
			a:  []interface{}{"hello world"},
			es: `"hello world"`,
		},
		{
			a:  []interface{}{"what happen\n"},
			es: "\"what happen\\n\"",
		},
		{
			a:  []interface{}{"what \\ happen"},
			es: "\"what \\\\ happen\"",
		},
		{
			a:  []interface{}{"{\"name\": \"test\"}"},
			es: "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			a:  []interface{}{"{\"name\": \"test\", \"age\": 18}"},
			es: "\"{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}\"",
		},
		{
			a:  []interface{}{"what \\ happen", 22, 12.4, true, false},
			es: "\"what \\\\ happen22 12.4 true false\"",
		},
	}
	b := NewBuffer(nil)
	for _, tt := range tests {
		b.Reset()
		b.WriteQuoteSprint(tt.a...)
		if tt.es != string(b.Bytes()) {
			t.Fatalf("get %s, want %s", string(b.Bytes()), tt.es)
		}
	}
}

func TestBuffer_WriteQuoteSprintf(t *testing.T) {
	tests := []struct {
		f  string
		a  []interface{}
		es string
	}{
		{
			f:  "%s",
			a:  []interface{}{"hello world"},
			es: `"hello world"`,
		},
		{
			f:  "%s",
			a:  []interface{}{"what happen\n"},
			es: "\"what happen\\n\"",
		},
		{
			a:  []interface{}{"what \\ happen"},
			es: "\"what \\\\ happen\"",
		},
		{
			f:  "{\"name\": \"test\"}",
			es: "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			f:  "%s %d %t %s",
			a:  []interface{}{"what \\ happen", 22, true, "{\"name\": \"test\"}"},
			es: "\"what \\\\ happen 22 true {\\\"name\\\": \\\"test\\\"}\"",
		},
	}

	b := NewBuffer(nil)
	for _, tt := range tests {
		b.Reset()
		b.WriteQuoteSprintf(tt.f, tt.a...)
		if tt.es != string(b.Bytes()) {
			t.Fatalf("get %s, want %s", string(b.Bytes()), tt.es)
		}
	}
}
