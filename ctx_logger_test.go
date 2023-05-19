package olog

import (
	"bytes"
	"context"
	"fmt"
	"testing"
)

func TestCtxPrintf(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		lv   Level
	}{
		{
			name: "Errorf",
			lv:   ERROR,
		},
		{
			name: "Warnf",
			lv:   WARN,
		},
		{
			name: "Noticef",
			lv:   NOTICE,
		},
		{
			name: "Infof",
			lv:   INFO,
		},
		{
			name: "Debugf",
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	logging := func(logger Logger, method string) {
		switch method {
		case "Errorf":
			logger.Errorf("test %s", "printf")
		case "Warnf":
			logger.Warnf("test %s", "printf")
		case "Noticef":
			logger.Noticef("test %s", "printf")
		case "Infof":
			logger.Infof("test %s", "printf")
		case "Debugf":
			logger.Debugf("test %s", "printf")
		}
	}

	for _, tt := range tests {
		SetEncode(PLAIN)
		logger := WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test printf")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		logger = WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test printf") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestCtxPrint(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		lv   Level
	}{
		{
			name: "Error",
			lv:   ERROR,
		},
		{
			name: "Warn",
			lv:   WARN,
		},
		{
			name: "Notice",
			lv:   NOTICE,
		},
		{
			name: "Info",
			lv:   INFO,
		},
		{
			name: "Debug",
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	logging := func(logger Logger, method string) {
		switch method {
		case "Error":
			logger.Error("test ", "print")
		case "Warn":
			logger.Warn("test ", "print")
		case "Notice":
			logger.Notice("test ", "print")
		case "Info":
			logger.Info("test ", "print")
		case "Debug":
			logger.Debug("test ", "print")
		}
	}

	for _, tt := range tests {
		SetEncode(PLAIN)
		logger := WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test print")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		logger = WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test print") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestCtxPrintw(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		lv   Level
	}{
		{
			name: "Errorw",
			lv:   ERROR,
		},
		{
			name: "Warnw",
			lv:   WARN,
		},
		{
			name: "Noticew",
			lv:   NOTICE,
		},
		{
			name: "Infow",
			lv:   INFO,
		},
		{
			name: "Debugw",
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))
	logging := func(logger Logger, method string) {
		switch method {
		case "Errorw":
			logger.Errorw("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		case "Warnw":
			logger.Warnw("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		case "Noticew":
			logger.Noticew("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		case "Infow":
			logger.Infow("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		case "Debugw":
			logger.Debugw("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		}
	}

	for _, tt := range tests {
		SetEncode(PLAIN)
		logger := WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want := fmt.Sprintf("\t%s\t%s\t%s=%s\t%s=%s\n", tt.lv.String(), "test", "age", "18", "addr", "new york")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		logger = WithContext(GetLogger(), context.Background())
		logging(logger, tt.name)
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s","%s":%d,"%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test", "age", 18, "addr", "new york") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestWithContext(t *testing.T) {
	setDefLogger(newLogger())

	ctx := context.WithValue(context.Background(), "uid", 3)
	SetDefCtxHandle(func(ctx context.Context) []Field {
		var fs []Field
		uid, ok := ctx.Value("uid").(int)
		if ok {
			fs = append(fs, Field{Key: "uid", Value: uid})
		}

		name, ok := ctx.Value("name").(string)
		if ok {
			fs = append(fs, Field{Key: "name", Value: name})
		}
		return fs
	})

	SetEncode(PLAIN)
	fields1 := []Field{{Key: "uid", Value: 3}}
	l := WithContext(GetLogger(), ctx)
	l.Debug("test 1")
	validateFields(t, l, fields1)

	fields2 := []Field{{Key: "name", Value: "bob"}, {Key: "uid", Value: 3}}
	ctx = context.WithValue(context.Background(), "name", "bob")
	l = WithContext(l, ctx)
	l.Debug("test 2")
	validateFields(t, l, fields2)

	fields3 := fields2
	l = WithContext(l, context.Background())
	l.Debug("test 3")
	validateFields(t, l, fields3)

	l = WithEntries(l, map[string]interface{}{
		"ip":      "127.0.0.1",
		"score":   99.9,
		"success": true,
	})
	fields4 := l.buildFields()
	l.Debug("test 4")
	if len(fields4) != 5 {
		t.Fatal("fields length not correct")
	}

	fields5 := []Field{{Key: "name", Value: "linda"}}
	fields5 = append(fields5, fields4...)
	l = WithContext(l, context.WithValue(context.Background(), "name", "linda"))
	l.Debug("test 5")
	validateFields(t, l, fields5)

	l.Log(Record{Level: DEBUG, Caller: Disable, LevelTag: "print", Stack: Enable, StackSize: 0,
		MsgOrFormat: "test 6"})
}

func validateFields(t *testing.T, l Logger, fields []Field) {
	fs := l.buildFields()
	if len(fs) != len(fields) {
		t.Fatal("fields length not equal")
	}

	for i, f := range fs {
		if f.Key != fields[i].Key {
			t.Fatalf("field key not equal, want = %s, got = %s", fields[i].Key, f.Key)
		}
		if f.Value != fields[i].Value {
			t.Fatalf("field value not equal, want = %v, got = %v", fields[i].Value, f.Value)
		}
	}
}
