package olog

// Logger is an interface that defines the methods for logging.
type Logger interface {
	// Log writes a log message with the given log level.
	Log(opt LogOption, args ...any)
	// Logf writes a log message with the given log level.
	Logf(opt LogOption, format string, args ...any)

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

	log(opt LogOption, args ...any)
	logf(opt LogOption, format string, args ...any)
	buildFields(fields ...Field) []Field
}

// Field is a struct that represents a key-value pair of additional data to include in a log message.
type Field struct {
	Key   string
	Value any
}

type EnableOp uint8

const (
	EnableDefault EnableOp = iota
	EnableOpen
	EnableClose
)

type LogOption struct {
	Level        Level
	EnableCaller EnableOp
	EnableStack  EnableOp
	StackSize    uint8
	CallerSkip   int8
	LevelTag     string
	Fields       []Field
}

// logOption is a struct that represents options to use when logging a message.
type logOption struct {
	level        Level   // level is the severity level of the log message.
	enableCaller bool    // enableCaller indicates whether to include caller information in the log message.
	enableColor  bool    // enableColor indicates whether to enable colorized output for the levelTag on plain encoding.
	enableStack  bool    // enableStack indicates whether to include stack trace information in the log message.
	stackSize    uint8   // stackSize is the maximum number of stack frames to include in the log message.
	callerSkip   int8    // callerSkip is the number of stack frames to skip to find the caller information.
	msgType      msgType // msgType is the type of the log message.
	msgArgs      []any   // msgArgs is a slice of arguments to the log message.
	msgOrFormat  string  // msgOrFormat is the format string of the log message.

	// tag is the string representation of the severity level
	// The default debug, info, warn, error, and fatal correspond to DEBUG, INFO, WARN, ERROR, and FATAL log levels respectively
	// users can also customize semantic tags, such as slow.
	tag        string
	timeFormat string
	fields     []Field // fields is a slice of key-value pairs of additional data to include in the log message.
}

// defCallerSkip is the default number of stack frames to skip to find the caller information.
const defCallerSkip = 5

// defStackSize is the default maximum number of stack frames to include in the log message.
const defStackSize = 5

type msgType int8

const (
	msgTypePrint msgType = iota + 1
	msgTypePrintf
	msgTypePrintMsg
)

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
