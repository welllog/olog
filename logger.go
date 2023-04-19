package olog

import (
	"os"
	"time"
)

// logger represents a logger instance with configurable options
type logger struct {
	level        Level                        // the minimum level of logging to output
	enableCaller bool                         // flag indicating whether to log the caller information
	enableColor  bool                         // flag indicating whether to use colorized output for levelTag on plain encoding
	timeFormat   string                       // time format to use for logging
	encode       func(o *logOption, w Writer) // encoding function to use for logging
	writer       Writer                       // writer to output log to
}

// NewLogger returns a new Logger instance with optional configurations
func NewLogger(opts ...LoggerOption) Logger {
	return newLogger(opts...)
}

// newLogger creates a new logger instance with default options that can be customized with the provided options
func newLogger(opts ...LoggerOption) *logger {
	l := logger{
		enableCaller: true,
		enableColor:  true,
		timeFormat:   time.RFC3339,
		encode:       jsonEncode,
		writer:       csWriter,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}

// LoggerOption is a functional option type for configuring a logger instance
type LoggerOption func(*logger)

// WithLoggerLevel sets the minimum logging level for the logger instance
func WithLoggerLevel(level Level) LoggerOption {
	return func(l *logger) {
		l.level = level
	}
}

// WithLoggerCaller sets whether to log the caller information in the output
func WithLoggerCaller(enable bool) LoggerOption {
	return func(l *logger) {
		l.enableCaller = enable
	}
}

// WithLoggerColor sets whether to use colorized output for levelTag on plain encoding, not supporting windows.
func WithLoggerColor(enable bool) LoggerOption {
	return func(l *logger) {
		l.enableColor = enable
	}
}

// WithLoggerTimeFormat sets the time format to use for logging
func WithLoggerTimeFormat(format string) LoggerOption {
	return func(l *logger) {
		l.timeFormat = format
	}
}

// WithLoggerEncode sets the encoding type to use for logging
func WithLoggerEncode(e EncodeType) LoggerOption {
	return func(l *logger) {
		if e == PLAIN {
			l.encode = plainEncode
		} else {
			l.encode = jsonEncode
		}
	}
}

// WithLoggerWriter sets the output writer for the logger instance
func WithLoggerWriter(w Writer) LoggerOption {
	return func(l *logger) {
		l.writer = w
	}
}

func (l *logger) Log(level Level, opts ...LogOption) {
	l.log(level, opts...)
}

func (l *logger) Fatal(a ...any) {
	l.fatal(a...)
}

func (l *logger) Fatalf(format string, a ...any) {
	l.fatalf(format, a...)
}

func (l *logger) Fatalw(msg string, fields ...Field) {
	l.fatalw(msg, fields...)
}

func (l *logger) Error(a ...any) {
	l.error(a...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.errorf(format, a...)
}

func (l *logger) Errorw(msg string, fields ...Field) {
	l.errorw(msg, fields...)
}

func (l *logger) Warn(a ...any) {
	l.warn(a...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.warnf(format, a...)
}

func (l *logger) Warnw(msg string, fields ...Field) {
	l.warnw(msg, fields...)
}

func (l *logger) Notice(a ...any) {
	l.notice(a...)
}

func (l *logger) Noticef(format string, a ...any) {
	l.noticef(format, a...)
}

func (l *logger) Noticew(msg string, fields ...Field) {
	l.noticew(msg, fields...)
}

func (l *logger) Info(a ...any) {
	l.info(a...)
}

func (l *logger) Infof(format string, a ...any) {
	l.infof(format, a...)
}

func (l *logger) Infow(msg string, fields ...Field) {
	l.infow(msg, fields...)
}

func (l *logger) Debug(a ...any) {
	l.debug(a...)
}

func (l *logger) Debugf(format string, a ...any) {
	l.debugf(format, a...)
}

func (l *logger) Debugw(msg string, fields ...Field) {
	l.debugw(msg, fields...)
}

func (l *logger) Trace(a ...any) {
	l.trace(a...)
}

func (l *logger) Tracef(format string, a ...any) {
	l.tracef(format, a...)
}

func (l *logger) Tracew(msg string, fields ...Field) {
	l.tracew(msg, fields...)
}

func (l *logger) IsEnabled(level Level) bool {
	return level >= l.level
}

func (l *logger) log(level Level, opts ...LogOption) {
	if l.IsEnabled(level) {
		o := logOption{
			level:        level,
			enableCaller: l.enableCaller,
			callerSkip:   defCallerSkip,
		}
		for _, opt := range opts {
			opt(&o)
		}
		l.output(&o)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *logger) fatal(a ...any) {
	l.output(&logOption{
		level:        FATAL,
		enableCaller: l.enableCaller,
		msgType:      msgTypePrint,
		msgArgs:      a,
	})
	os.Exit(1)
}

func (l *logger) fatalf(format string, a ...any) {
	l.output(&logOption{
		level:        FATAL,
		enableCaller: l.enableCaller,
		msgType:      msgTypePrintf,
		msgArgs:      a,
		msgOrFormat:  format,
	})
	os.Exit(1)
}

func (l *logger) fatalw(msg string, fields ...Field) {
	l.output(&logOption{
		level:        FATAL,
		enableCaller: l.enableCaller,
		msgType:      msgTypePrintMsg,
		msgOrFormat:  msg,
		fields:       fields,
	})
	os.Exit(1)
}

func (l *logger) error(a ...any) {
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) errorf(format string, a ...any) {
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) errorw(msg string, fields ...Field) {
	if l.IsEnabled(ERROR) {
		l.output(&logOption{
			level:        ERROR,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) warn(a ...any) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) warnf(format string, a ...any) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) warnw(msg string, fields ...Field) {
	if l.IsEnabled(WARN) {
		l.output(&logOption{
			level:        WARN,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) notice(a ...any) {
	if l.IsEnabled(NOTICE) {
		l.output(&logOption{
			level:        NOTICE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) noticef(format string, a ...any) {
	if l.IsEnabled(NOTICE) {
		l.output(&logOption{
			level:        NOTICE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) noticew(msg string, fields ...Field) {
	if l.IsEnabled(NOTICE) {
		l.output(&logOption{
			level:        NOTICE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) info(a ...any) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) infof(format string, a ...any) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) infow(msg string, fields ...Field) {
	if l.IsEnabled(INFO) {
		l.output(&logOption{
			level:        INFO,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) debug(a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) debugf(format string, a ...any) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) debugw(msg string, fields ...Field) {
	if l.IsEnabled(DEBUG) {
		l.output(&logOption{
			level:        DEBUG,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) trace(a ...any) {
	if l.IsEnabled(TRACE) {
		l.output(&logOption{
			level:        TRACE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrint,
			msgArgs:      a,
		})
	}
}

func (l *logger) tracef(format string, a ...any) {
	if l.IsEnabled(TRACE) {
		l.output(&logOption{
			level:        TRACE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintf,
			msgArgs:      a,
			msgOrFormat:  format,
		})
	}
}

func (l *logger) tracew(msg string, fields ...Field) {
	if l.IsEnabled(TRACE) {
		l.output(&logOption{
			level:        TRACE,
			enableCaller: l.enableCaller,
			msgType:      msgTypePrintMsg,
			msgOrFormat:  msg,
			fields:       fields,
		})
	}
}

func (l *logger) buildFields(fields ...Field) []Field {
	return fields
}

func (l *logger) output(o *logOption) {
	if o.tag == "" {
		o.tag = o.level.String()
	}
	if o.enableCaller {
		if o.callerSkip <= 0 {
			o.callerSkip = defCallerSkip
		}
		o.file, o.line = getCaller(o.callerSkip)
	}
	o.fields = l.filterFields(o.fields)
	o.enableColor = l.enableColor
	o.timeFormat = l.timeFormat
	l.encode(o, l.writer)
}

func (l *logger) filterFields(fields []Field) []Field {
	n := len(fields)
	if n == 0 {
		return fields
	}

	set := make(map[string]struct{}, n)
	var remain int
	for idx, field := range fields {
		_, ok := filterField[field.Key]
		if !ok {
			_, ok = set[field.Key]
			if !ok {
				set[field.Key] = struct{}{}
				fields[remain], fields[idx] = fields[idx], fields[remain]
				remain++
			}
		}
	}
	return fields[:remain]
}

func (l *logger) clone() *logger {
	return &logger{
		level:        l.level,
		enableCaller: l.enableCaller,
		enableColor:  l.enableColor,
		timeFormat:   l.timeFormat,
		encode:       l.encode,
		writer:       l.writer,
	}
}
