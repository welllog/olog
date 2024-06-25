package olog

import (
	"os"
	"sync"
	"time"

	"github.com/welllog/olog/encoder"
)

// defCallerSkip is the default number of stack frames to skip to find the caller information.
const defCallerSkip = 5

// defStackSize is the default maximum number of stack frames to include in the log message.
const defStackSize = 5

// Declare a new sync.Pool object, which allows for efficient
// re-use of objects across goroutines.
var bufPool = sync.Pool{
	New: func() interface{} {
		var b [256]byte
		return encoder.NewBuffer(b[:0])
	},
}

// getBuf to retrieve a *Buffer from the pool.
func getBuf() *encoder.Buffer {
	return bufPool.Get().(*encoder.Buffer)
}

// putBuf to return a *Buffer to the pool.
func putBuf(buf *encoder.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}

// BeforeEncHook is the hook function to execute before encoding the log message.
type BeforeEncHook func(string, []any) (string, []any)

// AfterEncHook is the hook function to execute after encoding the log message.
type AfterEncHook func([]byte) []byte

// logger represents a logger instance with configurable options
type logger struct {
	app       string          // the name of the application
	level     Level           // the minimum level of logging to output
	caller    EnableOp        // flag indicating whether to log the caller information
	color     EnableOp        // flag indicating whether to use colorized output for levelTag on plain encoding
	shortFile EnableOp        // flag indicating whether to use short file name in the log message
	encType   EncodeType      // the encoding type to use for encoding the log message
	timeFmt   string          // time format to use for logging
	enc       EncodeFunc      // enc to use for encoding the log message
	wr        Writer          // wr to output log to
	beforeEnc []BeforeEncHook // beforeEnc to execute before encoding the log message
	afterEnc  []AfterEncHook  // afterEnc to execute after encoding the log message
}

// NewLogger returns a new Logger instance with optional configurations
func NewLogger(opts ...LoggerOption) Logger {
	return newLogger(opts...)
}

// newLogger creates a new logger instance with default options that can be customized with the provided options
func newLogger(opts ...LoggerOption) *logger {
	l := logger{
		caller:    Enable,
		color:     Enable,
		shortFile: Enable,
		encType:   JSON,
		timeFmt:   time.RFC3339,
		wr:        csWriter,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}

// LoggerOption is a functional option type for configuring a logger instance
type LoggerOption func(*logger)

// WithLoggerAppName sets the name of the application
func WithLoggerAppName(name string) LoggerOption {
	return func(l *logger) {
		l.app = EscapedString(name)
	}
}

// WithLoggerLevel sets the minimum logging level for the logger instance
func WithLoggerLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = level
	}
}

// WithLoggerCaller sets whether to log the caller information in the output
func WithLoggerCaller(enable bool) LoggerOption {
	return func(l *logger) {
		if enable {
			l.caller = Enable
		} else {
			l.caller = Disable
		}
	}
}

// WithLoggerColor sets whether to use colorized output for levelTag on plain encoding, not supporting windows.
func WithLoggerColor(enable bool) LoggerOption {
	return func(l *logger) {
		if enable {
			l.color = Enable
		} else {
			l.color = Disable
		}
	}
}

// WithLoggerShortFile sets whether to use short file name in the log message
func WithLoggerShortFile(enable bool) LoggerOption {
	return func(l *logger) {
		if enable {
			l.shortFile = Enable
		} else {
			l.shortFile = Disable
		}
	}
}

// WithLoggerTimeFormat sets the time format to use for logging
func WithLoggerTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFmt = EscapedString(format)
	}
}

// WithLoggerEncode sets the encoding type to use for logging
func WithLoggerEncode(e EncodeType) LoggerOption {
	return func(l *logger) {
		switch e {
		case PLAIN, JSON:
			l.encType = e
		default:
			l.encType = JSON
		}
	}
}

// WithLoggerEncodeFunc sets the encoder to use for logging
func WithLoggerEncodeFunc(e EncodeFunc) LoggerOption {
	return func(l *logger) {
		l.enc = e
		l.encType = -1
	}
}

// WithLoggerWriter sets the output writer for the logger instance
func WithLoggerWriter(w Writer) LoggerOption {
	return func(l *logger) {
		l.wr = w
	}
}

// WithLoggerBeforeEnc adds a function to execute before encoding the log message
func WithLoggerBeforeEnc(f ...BeforeEncHook) LoggerOption {
	return func(l *logger) {
		l.beforeEnc = append(l.beforeEnc, f...)
	}
}

// WithLoggerAfterEnc adds a function to execute after encoding the log message
func WithLoggerAfterEnc(f ...AfterEncHook) LoggerOption {
	return func(l *logger) {
		l.afterEnc = append(l.afterEnc, f...)
	}
}

func (l *logger) Log(r Record) {
	if l.IsEnabled(r.Level) {
		l.log(r)
	}
}

func (l *logger) Fatal(args ...any) {
	l.fatal(args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.fatalf(format, args...)
}

func (l *logger) Fatalw(msg string, fields ...Field) {
	l.fatalw(msg, fields...)
}

func (l *logger) Error(args ...any) {
	l.error(args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.errorf(format, args...)
}

func (l *logger) Errorw(msg string, fields ...Field) {
	l.errorw(msg, fields...)
}

func (l *logger) Warn(args ...any) {
	l.warn(args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.warnf(format, args...)
}

func (l *logger) Warnw(msg string, fields ...Field) {
	l.warnw(msg, fields...)
}

func (l *logger) Notice(args ...any) {
	l.notice(args...)
}

func (l *logger) Noticef(format string, args ...any) {
	l.noticef(format, args...)
}

func (l *logger) Noticew(msg string, fields ...Field) {
	l.noticew(msg, fields...)
}

func (l *logger) Info(args ...any) {
	l.info(args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.infof(format, args...)
}

func (l *logger) Infow(msg string, fields ...Field) {
	l.infow(msg, fields...)
}

func (l *logger) Debug(args ...any) {
	l.debug(args...)
}

func (l *logger) Debugf(format string, args ...any) {
	l.debugf(format, args...)
}

func (l *logger) Debugw(msg string, fields ...Field) {
	l.debugw(msg, fields...)
}

func (l *logger) Trace(args ...any) {
	l.trace(args...)
}

func (l *logger) Tracef(format string, args ...any) {
	l.tracef(format, args...)
}

func (l *logger) Tracew(msg string, fields ...Field) {
	l.tracew(msg, fields...)
}

func (l *logger) IsEnabled(level Level) bool {
	return level >= l.level
}

func (l *logger) log(r Record) {
	if r.Stack == Default {
		r.Stack = Disable
	}

	if r.StackSize == 0 {
		r.StackSize = defStackSize
	}

	r.CallerSkip = defCallerSkip + r.CallerSkip

	l.output(r)
}

func (l *logger) fatal(a ...any) {
	l.output(Record{
		Level:   FATAL,
		MsgArgs: a,
		OsExit:  true,
	})
}

func (l *logger) fatalf(format string, a ...any) {
	l.output(Record{
		Level:       FATAL,
		MsgOrFormat: format,
		MsgArgs:     a,
		OsExit:      true,
	})
}

func (l *logger) fatalw(msg string, fields ...Field) {
	l.output(Record{
		Level:       FATAL,
		MsgOrFormat: msg,
		Fields:      fields,
		OsExit:      true,
	})
}

func (l *logger) error(a ...any) {
	if l.IsEnabled(ERROR) {
		l.output(Record{
			Level:   ERROR,
			MsgArgs: a,
		})
	}
}

func (l *logger) errorf(format string, a ...any) {
	if l.IsEnabled(ERROR) {
		l.output(Record{
			Level:       ERROR,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) errorw(msg string, fields ...Field) {
	if l.IsEnabled(ERROR) {
		l.output(Record{
			Level:       ERROR,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) warn(a ...any) {
	if l.IsEnabled(WARN) {
		l.output(Record{
			Level:   WARN,
			MsgArgs: a,
		})
	}
}

func (l *logger) warnf(format string, a ...any) {
	if l.IsEnabled(WARN) {
		l.output(Record{
			Level:       WARN,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) warnw(msg string, fields ...Field) {
	if l.IsEnabled(WARN) {
		l.output(Record{
			Level:       WARN,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) notice(a ...any) {
	if l.IsEnabled(NOTICE) {
		l.output(Record{
			Level:   NOTICE,
			MsgArgs: a,
		})
	}
}

func (l *logger) noticef(format string, a ...any) {
	if l.IsEnabled(NOTICE) {
		l.output(Record{
			Level:       NOTICE,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) noticew(msg string, fields ...Field) {
	if l.IsEnabled(NOTICE) {
		l.output(Record{
			Level:       NOTICE,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) info(a ...any) {
	if l.IsEnabled(INFO) {
		l.output(Record{
			Level:   INFO,
			MsgArgs: a,
		})
	}
}

func (l *logger) infof(format string, a ...any) {
	if l.IsEnabled(INFO) {
		l.output(Record{
			Level:       INFO,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) infow(msg string, fields ...Field) {
	if l.IsEnabled(INFO) {
		l.output(Record{
			Level:       INFO,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) debug(a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(Record{
			Level:   DEBUG,
			MsgArgs: a,
		})
	}
}

func (l *logger) debugf(format string, a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(Record{
			Level:       DEBUG,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) debugw(msg string, fields ...Field) {
	if l.IsEnabled(DEBUG) {
		l.output(Record{
			Level:       DEBUG,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) trace(a ...any) {
	if l.IsEnabled(TRACE) {
		l.output(Record{
			Level:     TRACE,
			Stack:     Enable,
			StackSize: defStackSize,
			MsgArgs:   a,
		})
	}
}

func (l *logger) tracef(format string, a ...any) {
	if l.IsEnabled(TRACE) {
		l.output(Record{
			Level:       TRACE,
			Stack:       Enable,
			StackSize:   defStackSize,
			MsgOrFormat: format,
			MsgArgs:     a,
		})
	}
}

func (l *logger) tracew(msg string, fields ...Field) {
	if l.IsEnabled(TRACE) {
		l.output(Record{
			Level:       TRACE,
			Stack:       Enable,
			StackSize:   defStackSize,
			MsgOrFormat: msg,
			Fields:      fields,
		})
	}
}

func (l *logger) buildFields(fields ...Field) []Field {
	return fields
}

func (l *logger) output(r Record) {
	if r.Caller == Default {
		r.Caller = l.caller
	}

	if r.ShortFile == Default {
		r.ShortFile = l.shortFile
	}

	if r.CallerSkip <= 0 {
		r.CallerSkip = defCallerSkip
	}

	if r.LevelTag == "" {
		r.LevelTag = r.Level.String()
	} else {
		r.LevelTag = EscapedString(r.LevelTag)
	}

	if r.Time.IsZero() {
		r.Time = time.Now()
	}

	r.App = l.app
	r.TimeFmt = l.timeFmt

	for _, f := range l.beforeEnc {
		r.MsgOrFormat, r.MsgArgs = f(r.MsgOrFormat, r.MsgArgs)
	}

	buf := getBuf()

	switch l.encType {
	case PLAIN:
		plainEncode(r, buf, l.color.IsOpen())
	case -1:
		l.enc(r, buf)
	default:
		jsonEncode(r, buf)
	}

	data := buf.Bytes()
	for _, f := range l.afterEnc {
		data = f(data)
	}
	_, _ = l.wr.Write(r.Level, data)

	putBuf(buf)

	if r.OsExit {
		os.Exit(1)
	}
}

func (l *logger) clone() *logger {
	return &logger{
		app:       l.app,
		level:     l.level,
		caller:    l.caller,
		color:     l.color,
		shortFile: l.shortFile,
		encType:   l.encType,
		timeFmt:   l.timeFmt,
		enc:       l.enc,
		wr:        l.wr,
		afterEnc:  l.afterEnc,
		beforeEnc: l.beforeEnc,
	}
}
