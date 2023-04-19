package olog

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"
)

func initTestLogger() {
	SetLevel(TRACE)
	SetTimeFormat("")
	SetColor(false)
	SetCaller(false)
	SetEncode(JSON)
	SetWriter(csWriter)
}

func TestPrintf(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(string, ...any)
		lv   Level
	}{
		{
			name: "Errorf",
			fn:   Errorf,
			lv:   ERROR,
		},
		{
			name: "Warnf",
			fn:   Warnf,
			lv:   WARN,
		},
		{
			name: "Noticef",
			fn:   Noticef,
			lv:   NOTICE,
		},
		{
			name: "Infof",
			fn:   Infof,
			lv:   INFO,
		},
		{
			name: "Debugf",
			fn:   Debugf,
			lv:   DEBUG,
		},
		{
			name: "Tracef",
			fn:   Tracef,
			lv:   TRACE,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test %s", "printf")
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test printf")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test %s", "printf")
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test printf") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestPrint(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(...any)
		lv   Level
	}{
		{
			name: "Error",
			fn:   Error,
			lv:   ERROR,
		},
		{
			name: "Warn",
			fn:   Warn,
			lv:   WARN,
		},
		{
			name: "Notice",
			fn:   Notice,
			lv:   NOTICE,
		},
		{
			name: "Info",
			fn:   Info,
			lv:   INFO,
		},
		{
			name: "Debug",
			fn:   Debug,
			lv:   DEBUG,
		},
		{
			name: "Trace",
			fn:   Trace,
			lv:   TRACE,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test ", "print")
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test print")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test ", "print")
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test print") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestPrintw(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(string, ...Field)
		lv   Level
	}{
		{
			name: "Errorw",
			fn:   Errorw,
			lv:   ERROR,
		},
		{
			name: "Warnw",
			fn:   Warnw,
			lv:   WARN,
		},
		{
			name: "Noticew",
			fn:   Noticew,
			lv:   NOTICE,
		},
		{
			name: "Infow",
			fn:   Infow,
			lv:   INFO,
		},
		{
			name: "Debugw",
			fn:   Debugw,
			lv:   DEBUG,
		},
		{
			name: "Tracew",
			fn:   Tracew,
			lv:   TRACE,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		want := fmt.Sprintf("\t%s\t%s\t%s=%s\t%s=%s\n", tt.lv.String(), "test", "age", "18", "addr", "new york")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s","%s":%d,"%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test", "age", 18, "addr", "new york") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestSetLevel(t *testing.T) {
	initTestLogger()

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	logging := func() {
		Error("test")
		Warn("test")
		Notice("test")
		Info("test")
		Debug("test")
		Trace("test")
	}

	getLines := func() int {
		var count int
		for {
			line, _ := buf.ReadBytes('\n')
			if len(line) == 0 {
				break
			}
			count++
		}
		return count
	}

	SetLevel(ERROR)
	logging()
	lines := getLines()
	want := 1
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(WARN)
	logging()
	lines = getLines()
	want = 2
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(NOTICE)
	logging()
	lines = getLines()
	want = 3
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(INFO)
	logging()
	lines = getLines()
	want = 4
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(DEBUG)
	logging()
	lines = getLines()
	want = 5
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(TRACE)
	logging()
	lines = getLines()
	want = 6
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}
}

func TestLogFacade(t *testing.T) {
	initTestLogger()
	SetColor(true)
	SetCaller(true)

	tests := []struct {
		arr    []any
		f      string
		a      []any
		msg    string
		fields []Field
	}{
		{
			arr: []any{"t1", "t2"},
			f:   "t%d%s",
			a:   []any{3, "t4"},
			msg: "t5",
			fields: []Field{
				{Key: "name", Value: "bob"},
				{Key: "age", Value: 18},
			},
		},
	}

	for _, tt := range tests {
		logging(tt)
	}

	SetEncode(PLAIN)
	for _, tt := range tests {
		logging(tt)
	}

	SetTimeFormat("2006/01/02 15/04/05")
	for _, tt := range tests {
		logging(tt)
	}

	f, err := os.Create("test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	SetWriter(NewWriter(f))
	SetTimeFormat("")
	for _, tt := range tests {
		logging(tt)
	}

	SetColor(false)
	for _, tt := range tests {
		logging(tt)
	}
}

func TestLog(t *testing.T) {
	initTestLogger()
	SetCaller(true)
	SetEncode(PLAIN)
	Log(INFO, WithPrintMsg("test log"), WithTag("stat"), WithFields(Field{Key: "name", Value: "bob"}))
	Log(WARN, WithPrintMsg("test log"), WithTag("slow"), WithCaller(false))
	Log(WARN, WithPrintMsg("test log"), WithTag("slow"), WithCallerSkip(1), WithCallerSkip(-1))
}

type customLogger struct {
	Logger
}

func (l *customLogger) Slow(a ...any) {
	l.Log(WARN, WithPrint(a...), WithTag("slow"), WithCallerSkipOne)
}

func (l *customLogger) Stat(a ...any) {
	Log(INFO, WithPrint(a...), WithTag("stat"), WithCallerSkipOne)
}

func TestWrapLogger(t *testing.T) {
	initTestLogger()
	SetCaller(true)
	SetEncode(PLAIN)
	l := customLogger{
		Logger: GetLogger(),
	}
	l.Slow("test slow")
	l.Stat("test stat")
}

func logging(tt struct {
	arr    []any
	f      string
	a      []any
	msg    string
	fields []Field
}) {
	Debug(tt.arr...)

	Infof(tt.f, tt.a...)

	Warnw(tt.msg, tt.fields...)
}

func TestPlainOutput(t *testing.T) {
	SetEncode(PLAIN)
	Trace("hello world")
	Tracew("hello", Field{Key: "name", Value: "bob"})
	Debug("hello world")
	Debugw("hello", Field{Key: "name", Value: "bob"})
	Info("hello world")
	Infow("hello", Field{Key: "name", Value: "linda"}, Field{Key: "age", Value: 18})
	Notice("hello world")
	Noticef("hello %s", "world")
	Warnf("hello %s", "world")
	Warnw("hello", Field{Key: "order_no", Value: "AWESDDF"})
	Error("hello world")
	Errorw("hello world", Field{Key: "success", Value: true})
	Log(DEBUG, WithTag("start"), WithPrintMsg("hello world"), WithCaller(false),
		WithFields(Field{Key: "price", Value: 32.5}))
}

func BenchmarkInfo(b *testing.B) {
	b.Run("std.logger", func(b *testing.B) {
		logger := log.New(discard{}, "", log.Ldate|log.Ltime|log.Lmsgprefix)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println("test message", "name", "bob", "age", 18, "success", true)
			}
		})
	})

	b.Run("olog.json", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerCaller(false))

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("test message", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
			}
		})
	})

	b.Run("olog.plain", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerCaller(false), WithLoggerEncode(PLAIN))

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("test message", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
			}
		})
	})

	b.Run("olog.ctx.json", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerCaller(false))
		logger = WithContext(logger, context.Background())
		logger = WithContext(logger, context.Background())

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("test message", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
			}
		})
	})

	b.Run("olog.ctx.plain", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerCaller(false), WithLoggerEncode(PLAIN))
		logger = WithContext(logger, context.Background())
		logger = WithContext(logger, context.Background())

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow("test message", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
			}
		})
	})
}

func BenchmarkWithCallerSimpleInfo(b *testing.B) {
	b.Run("std.logger", func(b *testing.B) {
		logger := log.New(discard{}, "", log.Ldate|log.Ltime|log.Lshortfile)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Print("test")
			logger.Printf("test %d", 2)
			logger.Println("test", "name", "bob", "age", 18, "success", true)
		}
	})

	b.Run("olog.json", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})))

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info("test")
			logger.Infof("test %d", 2)
			logger.Infow("test", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
		}
	})

	b.Run("olog.plain", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerEncode(PLAIN))

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info("test")
			logger.Infof("test %d", 2)
			logger.Infow("test", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
		}
	})

	b.Run("olog.ctx.json", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})))
		logger = WithContext(logger, context.Background())
		logger = WithContext(logger, context.Background())

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info("test")
			logger.Infof("test %d", 2)
			logger.Infow("test", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
		}
	})

	b.Run("olog.ctx.plain", func(b *testing.B) {
		logger := NewLogger(WithLoggerWriter(NewWriter(discard{})), WithLoggerEncode(PLAIN))
		logger = WithContext(logger, context.Background())
		logger = WithContext(logger, context.Background())

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info("test")
			logger.Infof("test %d", 2)
			logger.Infow("test", Field{Key: "name", Value: "bob"}, Field{Key: "age", Value: 18}, Field{Key: "success", Value: true})
		}
	})
}

type discard struct{}

func (d discard) Write(p []byte) (int, error) {
	return len(p), nil
}
