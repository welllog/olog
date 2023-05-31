package encoder

import (
	"testing"
)

func TestEPrint(t *testing.T) {
	tests := []struct {
		a  []any
		es string
	}{
		{
			a:  []any{"hello world"},
			es: `"hello world"`,
		},
		{
			a:  []any{"what happen\n"},
			es: "\"what happen\\n\"",
		},
		{
			a:  []any{"what \\ happen"},
			es: "\"what \\\\ happen\"",
		},
		{
			a:  []any{"{\"name\": \"test\"}"},
			es: "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			a:  []any{"{\"name\": \"test\", \"age\": 18}"},
			es: "\"{\\\"name\\\": \\\"test\\\", \\\"age\\\": 18}\"",
		},
		{
			a:  []any{"what \\ happen", 22, 12.4, true, false},
			es: "\"what \\\\ happen22 12.4 true false\"",
		},
	}

	e := JsonEncoder{&Buffer{}}
	for _, tt := range tests {
		e.Reset()
		e.WriteQuote()
		EPrint(e, tt.a...)
		e.WriteQuote()
		if tt.es != string(e.Bytes()) {
			t.Fatalf("get %s, want %s", string(e.Bytes()), tt.es)
		}
	}
}

func TestEPrintf(t *testing.T) {
	tests := []struct {
		f  string
		a  []any
		es string
	}{
		{
			f:  "%s",
			a:  []any{"hello world"},
			es: `"hello world"`,
		},
		{
			f:  "%s",
			a:  []any{"what happen\n"},
			es: "\"what happen\\n\"",
		},
		{
			a:  []any{"what \\ happen"},
			es: "\"what \\\\ happen\"",
		},
		{
			f:  "{\"name\": \"test\"}",
			es: "\"{\\\"name\\\": \\\"test\\\"}\"",
		},
		{
			f:  "%s %d %t %s",
			a:  []any{"what \\ happen", 22, true, "{\"name\": \"test\"}"},
			es: "\"what \\\\ happen 22 true {\\\"name\\\": \\\"test\\\"}\"",
		},
	}

	e := JsonEncoder{&Buffer{}}
	for _, tt := range tests {
		e.Reset()
		e.WriteQuote()
		EPrintf(e, tt.f, tt.a...)
		e.WriteQuote()
		if tt.es != string(e.Bytes()) {
			t.Fatalf("get %s, want %s", string(e.Bytes()), tt.es)
		}
	}
}

func TestEPrint1(t *testing.T) {
	buf := PlainEncoder{&Buffer{}}
	n, err := EPrint(buf, "hello", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 10 {
		t.Errorf("unexpected number of bytes written: %d", n)
	}
	if string(buf.Bytes()) != "helloworld" {
		t.Errorf("unexpected output: %q", string(buf.Bytes()))
	}
}

func TestEPrintf1(t *testing.T) {
	buf := PlainEncoder{&Buffer{}}
	n, err := EPrintf(buf, "hello %s", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 11 {
		t.Errorf("unexpected number of bytes written: %d", n)
	}
	if string(buf.Bytes()) != "hello world" {
		t.Errorf("unexpected output: %q", string(buf.Bytes()))
	}
}

func TestEPrintfNoArgs(t *testing.T) {
	buf := PlainEncoder{&Buffer{}}
	n, err := EPrintf(buf, "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("unexpected number of bytes written: %d", n)
	}
	if string(buf.Bytes()) != "hello" {
		t.Errorf("unexpected output: %q", string(buf.Bytes()))
	}
}

func TestEPrintfEmptyFormat(t *testing.T) {
	buf := PlainEncoder{&Buffer{}}
	n, err := EPrintf(buf, "", "hello", "world")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 10 {
		t.Errorf("unexpected number of bytes written: %d", n)
	}
	if string(buf.Bytes()) != "helloworld" {
		t.Errorf("unexpected output: %q", string(buf.Bytes()))
	}
}

func TestEPrintfNoArgsEmptyFormat(t *testing.T) {
	buf := PlainEncoder{&Buffer{}}
	n, err := EPrintf(buf, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("unexpected number of bytes written: %d", n)
	}
	if string(buf.Bytes()) != "" {
		t.Errorf("unexpected output: %q", string(buf.Bytes()))
	}
}
