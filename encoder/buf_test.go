package encoder

import (
	"bytes"
	"testing"
	"time"
)

func TestBuffer_Write(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "empty buffer",
			data: []byte{},
			want: []byte{},
		},
		{
			name: "non-empty buffer",
			data: []byte("hello world"),
			want: []byte("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			_, err := buf.Write(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteString(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []byte
	}{
		{
			name: "empty string",
			data: "",
			want: []byte{},
		},
		{
			name: "non-empty string",
			data: "hello world",
			want: []byte("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			_, err := buf.WriteString(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteByte(t *testing.T) {
	tests := []struct {
		name string
		data byte
		want []byte
	}{
		{
			name: "zero byte",
			data: 0,
			want: []byte{0},
		},
		{
			name: "non-zero byte",
			data: 'a',
			want: []byte{'a'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			err := buf.WriteByte(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteRune(t *testing.T) {
	tests := []struct {
		name string
		data rune
		want []byte
	}{
		{
			name: "ASCII character",
			data: 'a',
			want: []byte{'a'},
		},
		{
			name: "non-ASCII character",
			data: '世',
			want: []byte{0xe4, 0xb8, 0x96},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			_, err := buf.WriteRune(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteTime(t *testing.T) {
	tests := []struct {
		name   string
		layout string
		time   time.Time
		want   []byte
	}{
		{
			name:   "RFC3339 format",
			layout: time.RFC3339,
			time:   time.Date(2021, 8, 1, 12, 0, 0, 0, time.UTC),
			want:   []byte("2021-08-01T12:00:00Z"),
		},
		{
			name:   "custom format",
			layout: "2006-01-02 15:04:05",
			time:   time.Date(2021, 8, 1, 12, 0, 0, 0, time.UTC),
			want:   []byte("2021-08-01 12:00:00"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteTime(tt.time, tt.layout)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteInt64(t *testing.T) {
	tests := []struct {
		name string
		data int64
		want []byte
	}{
		{
			name: "zero",
			data: 0,
			want: []byte("0"),
		},
		{
			name: "positive number",
			data: 12345,
			want: []byte("12345"),
		},
		{
			name: "negative number",
			data: -12345,
			want: []byte("-12345"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteInt64(tt.data)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteUint64(t *testing.T) {
	tests := []struct {
		name string
		data uint64
		want []byte
	}{
		{
			name: "zero",
			data: 0,
			want: []byte("0"),
		},
		{
			name: "positive number",
			data: 12345,
			want: []byte("12345"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteUint64(tt.data)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteFloat(t *testing.T) {
	tests := []struct {
		name    string
		data    float64
		fmt     byte
		bitSize int
		want    []byte
	}{
		{
			name:    "zero",
			data:    0,
			fmt:     'f',
			bitSize: 64,
			want:    []byte("0"),
		},
		{
			name:    "positive number",
			data:    12.345,
			fmt:     'f',
			bitSize: 64,
			want:    []byte("12.345"),
		},
		{
			name:    "negative number",
			data:    -12.345,
			fmt:     'f',
			bitSize: 64,
			want:    []byte("-12.345"),
		},
		{
			name:    "exponential notation",
			data:    12345,
			fmt:     'e',
			bitSize: 64,
			want:    []byte("1.2345e+04"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteFloat(tt.data, tt.fmt, tt.bitSize)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteBool(t *testing.T) {
	tests := []struct {
		name string
		data bool
		want []byte
	}{
		{
			name: "true",
			data: true,
			want: []byte("true"),
		},
		{
			name: "false",
			data: false,
			want: []byte("false"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteBool(tt.data)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}

func TestBuffer_WriteBase64(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "empty slice",
			data: []byte{},
			want: []byte{},
		},
		{
			name: "single byte",
			data: []byte{0x01},
			want: []byte("AQ=="),
		},
		{
			name: "multiple bytes",
			data: []byte{0x01, 0x02, 0x03},
			want: []byte("AQID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &Buffer{}
			buf.WriteBase64(tt.data)
			if !bytes.Equal(buf.Bytes(), tt.want) {
				t.Errorf("got %q, want %q", buf.Bytes(), tt.want)
			}
		})
	}
}
