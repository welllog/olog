package olog

// Logger is an interface that defines the methods for logging.
type Logger interface {
	// Log writes a log message with the given level and options.
	// Users can define their own log methods according to this method.
	Log(level Level, opts ...LogOption)

	// Fatal writes a log message with the FATAL log level and call os.Exit(1).
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Fatalw(msg string, fields ...Field)

	// Error writes a log message with the ERROR log level.
	Error(args ...any)
	Errorf(format string, args ...any)
	Errorw(msg string, fields ...Field)

	// Warn writes a log message with the WARN log level.
	Warn(args ...any)
	Warnf(format string, args ...any)
	Warnw(msg string, fields ...Field)

	// Notice writes a log message with the NOTICE log level.
	Notice(args ...any)
	Noticef(format string, args ...any)
	Noticew(msg string, fields ...Field)

	// Info writes a log message with the INFO log level.
	Info(args ...any)
	Infof(format string, args ...any)
	Infow(msg string, fields ...Field)

	// Debug writes a log message with the DEBUG log level.
	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugw(msg string, fields ...Field)

	// Trace writes a log message with the TRACE log level.
	Trace(args ...any)
	Tracef(format string, args ...any)
	Tracew(msg string, fields ...Field)

	// IsEnabled returns whether the given log level is enabled or not.
	IsEnabled(level Level) bool

	log(level Level, opts ...LogOption)
	buildFields(fields ...Field) []Field
}

// Field is a struct that represents a key-value pair of additional data to include in a log message.
type Field struct {
	Key   string
	Value any
}

// LogOption is a function that modifies a logOption struct.
type LogOption func(*logOption)

// logOption is a struct that represents options to use when logging a message.
type logOption struct {
	level        Level   // level is the severity level of the log message.
	enableCaller bool    // enableCaller indicates whether to include caller information in the log message.
	enableColor  bool    // enableColor indicates whether to enable colorized output for the levelTag on plain encoding.
	msgType      msgType // msgType is the type of the log message.
	msgArgs      []any   // msgArgs is a slice of arguments to the log message.
	msgOrFormat  string  // msgOrFormat is the format string of the log message.
	callerSkip   int     // callerSkip is the number of stack frames to skip to find the caller information.
	file         string  // file is the file name of the log message.
	line         int     // line is the line number of the log message.

	// tag is the string representation of the severity level
	// The default debug, info, warn, error, and fatal correspond to DEBUG, INFO, WARN, ERROR, and FATAL log levels respectively
	// users can also customize semantic tags, such as slow.
	tag        string
	timeFormat string
	fields     []Field // fields is a slice of key-value pairs of additional data to include in the log message.
}

// defCallerSkip is the default number of stack frames to skip to find the caller information.
const defCallerSkip = 4

type msgType int8

const (
	msgTypePrint msgType = iota + 1
	msgTypePrintf
	msgTypePrintMsg
)

func WithTag(tag string) LogOption {
	return func(o *logOption) {
		o.tag = tag
	}
}

// WithCaller returns a LogOption that enables or disables logging the caller information.
func WithCaller(enable bool) LogOption {
	return func(o *logOption) {
		o.enableCaller = enable
	}
}

// WithCallerSkip returns a LogOption that sets the number of stack frames to skip when logging caller information.
func WithCallerSkip(skip int) LogOption {
	return func(o *logOption) {
		o.callerSkip += skip
	}
}

// WithFields returns a LogOption that appends additional key-value pairs to the logged message.
func WithFields(fields ...Field) LogOption {
	return func(o *logOption) {
		o.fields = append(o.fields, fields...)
	}
}

// WithPrint returns a LogOption that sets the type of the log message to print.
func WithPrint(args ...any) LogOption {
	return func(o *logOption) {
		if o.msgType > 0 {
			return
		}
		o.msgType = msgTypePrint
		o.msgArgs = args
	}
}

// WithPrintf returns a LogOption that sets the type of the log message to print.
func WithPrintf(format string, args ...any) LogOption {
	return func(o *logOption) {
		if o.msgType > 0 {
			return
		}
		o.msgType = msgTypePrintf
		o.msgOrFormat = format
		o.msgArgs = args
	}
}

// WithPrintMsg returns a LogOption that sets the type of the log message to print.
func WithPrintMsg(msg string) LogOption {
	return func(o *logOption) {
		if o.msgType > 0 {
			return
		}
		o.msgType = msgTypePrintMsg
		o.msgOrFormat = msg
	}
}

// WithCallerSkipOne is a LogOption that increments the number of stack frames to skip by 1 when logging caller information.
func WithCallerSkipOne(o *logOption) {
	o.callerSkip++
}

// WithCallerSkipTwo is a LogOption that increments the number of stack frames to skip by 2 when logging caller information.
func WithCallerSkipTwo(o *logOption) {
	o.callerSkip += 2
}

var (
	fieldTime    = "@timestamp"
	fieldLevel   = "level"
	fieldContent = "content"
	fieldCaller  = "caller"

	filterField = map[string]struct{}{
		fieldTime:    {},
		fieldLevel:   {},
		fieldContent: {},
		fieldCaller:  {},
	}
)
