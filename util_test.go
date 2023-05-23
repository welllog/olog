package olog

import "testing"

func TestEscapedString(t *testing.T) {
	tests := []struct {
		s  string
		es string
	}{
		{
			s:  "hello world",
			es: "hello world",
		},
		{
			s:  "what happen\n",
			es: "what happen\\n",
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
	for _, tt := range tests {
		if tt.es != EscapedString(tt.s) {
			t.Fatalf("get %q, want %q", EscapedString(tt.s), tt.es)
		}
	}
}
