package olog

import (
	"runtime"
	"time"

	"github.com/welllog/olog/encoder"
)

// Logger is an interface that defines the methods for logging.
type Logger interface {
	// Log writes a log message with the given log level.
	Log(r Record)

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

	log(r Record)
	buildFields(fields ...Field) []Field
}

// Field is a struct that represents a key-value pair of additional data to include in a log message.
type Field struct {
	Key   string
	Value any
}

// EnableOp is the enable of caller,stack,color information in the log message.
type EnableOp uint8

func (e EnableOp) IsOpen() bool {
	return e == Enable
}

const (
	Default EnableOp = iota
	Enable
	Disable
)

// EncodeFunc is a function that encodes a log message to a byte slice.
type EncodeFunc func(Record, *encoder.Buffer)

type Record struct {
	Level       Level    // Level is the severity level of the log message.
	Caller      EnableOp // Caller is the enable of caller information in the log message.
	Stack       EnableOp // Stack is the enable of stack trace information in the log message.
	ShortFile   EnableOp // ShortFile is the enable of short file name in the log message.
	StackSize   uint8    // StackSize is the maximum number of stack frames to include in the log message.
	CallerSkip  int8     // CallerSkip is the number of stack frames to skip to find the caller information.
	OsExit      bool     // OsExit is the enable of os.Exit(1) in the log message.
	MsgOrFormat string   // MsgOrFormat is the string representation of the log message
	MsgArgs     []any    // MsgArgs is the arguments of the log message
	Fields      []Field  // Fields is a slice of key-value pairs of additional data to include in the log message.
	LevelTag    string   // LevelTag is the string representation of the severity level
	App         string   // App is the name of the application that created the log message.
	TimeFmt     string   // TimeFmt is the format string of the log message.
	Time        time.Time
}

func (r Record) Frames() *runtime.Frames {
	var stackSize int
	if r.Stack.IsOpen() && r.StackSize > 0 {
		stackSize = int(r.StackSize)
	} else if r.Caller.IsOpen() {
		stackSize = 1
	} else {
		return nil
	}

	pc := make([]uintptr, stackSize)
	n := runtime.Callers(int(r.CallerSkip+1), pc)
	return runtime.CallersFrames(pc[:n])
}
