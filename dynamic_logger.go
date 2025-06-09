package olog

// DynamicLogger is a logger use global default logger, it has most same behavior as the default logger.
// But can set own application name, caller related options.
type DynamicLogger struct {
	Caller     EnableOp // Caller is the enable of caller information in the log message.
	ShortFile  EnableOp // ShortFile is the enable of short file name in the log message.
	CallerSkip int8
	app        string // app is the name of the application.
}

// SetAppName sets the name of the application for the DynamicLogger.
// If the application name is not set, it will use the default application name.
// The method is not thread-safe, so it should be called before any logging operations.
func (d *DynamicLogger) SetAppName(name string) {
	d.app = EscapedString(name)
}

func (d DynamicLogger) Log(r Record) {
	l := getDefLogger()
	if l.IsEnabled(r.Level) {
		if r.Caller == Default {
			r.Caller = d.Caller
		}

		if r.ShortFile == Default {
			r.ShortFile = d.ShortFile
		}

		if r.StackSize == 0 {
			r.StackSize = defStackSize
		}

		r.CallerSkip = defCallerSkip - 1 + d.CallerSkip + r.CallerSkip
		r.App = d.app

		l.output(r)
	}
}

func (d DynamicLogger) Fatal(args ...any) {
	d.printRecord(Record{
		Level:   FATAL,
		MsgArgs: args,
		OsExit:  true,
	})
}

func (d DynamicLogger) Fatalf(format string, args ...any) {
	d.printRecord(Record{
		Level:       FATAL,
		MsgOrFormat: format,
		MsgArgs:     args,
		OsExit:      true,
	})
}

func (d DynamicLogger) Fatalw(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       FATAL,
		MsgOrFormat: msg,
		Fields:      fields,
		OsExit:      true,
	})
}

func (d DynamicLogger) Error(args ...any) {
	d.printRecord(Record{
		Level:   ERROR,
		MsgArgs: args,
	})
}

func (d DynamicLogger) Errorf(format string, args ...any) {
	d.printRecord(Record{
		Level:       ERROR,
		MsgOrFormat: format,
		MsgArgs:     args,
	})
}

func (d DynamicLogger) Errorw(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       ERROR,
		MsgOrFormat: msg,
		Fields:      fields,
	})
}

func (d DynamicLogger) Warn(args ...any) {
	d.printRecord(Record{
		Level:   WARN,
		MsgArgs: args,
	})
}

func (d DynamicLogger) Warnf(format string, args ...any) {
	d.printRecord(Record{
		Level:       WARN,
		MsgOrFormat: format,
		MsgArgs:     args,
	})
}

func (d DynamicLogger) Warnw(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       WARN,
		MsgOrFormat: msg,
		Fields:      fields,
	})
}

func (d DynamicLogger) Notice(args ...any) {
	d.printRecord(Record{
		Level:   NOTICE,
		MsgArgs: args,
	})
}

func (d DynamicLogger) Noticef(format string, args ...any) {
	d.printRecord(Record{
		Level:       NOTICE,
		MsgOrFormat: format,
		MsgArgs:     args,
	})
}

func (d DynamicLogger) Noticew(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       NOTICE,
		MsgOrFormat: msg,
		Fields:      fields,
	})
}

func (d DynamicLogger) Info(args ...any) {
	d.printRecord(Record{
		Level:   INFO,
		MsgArgs: args,
	})
}

func (d DynamicLogger) Infof(format string, args ...any) {
	d.printRecord(Record{
		Level:       INFO,
		MsgOrFormat: format,
		MsgArgs:     args,
	})
}

func (d DynamicLogger) Infow(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       INFO,
		MsgOrFormat: msg,
		Fields:      fields,
	})
}

func (d DynamicLogger) Debug(args ...any) {
	d.printRecord(Record{
		Level:   DEBUG,
		MsgArgs: args,
	})
}

func (d DynamicLogger) Debugf(format string, args ...any) {
	d.printRecord(Record{
		Level:       DEBUG,
		MsgOrFormat: format,
		MsgArgs:     args,
	})
}

func (d DynamicLogger) Debugw(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       DEBUG,
		MsgOrFormat: msg,
		Fields:      fields,
	})
}

func (d DynamicLogger) Trace(args ...any) {
	d.printRecord(Record{
		Level:     TRACE,
		MsgArgs:   args,
		Stack:     Enable,
		StackSize: defStackSize,
	})
}

func (d DynamicLogger) Tracef(format string, args ...any) {
	d.printRecord(Record{
		Level:       TRACE,
		MsgOrFormat: format,
		MsgArgs:     args,
		Stack:       Enable,
		StackSize:   defStackSize,
	})
}

func (d DynamicLogger) Tracew(msg string, fields ...Field) {
	d.printRecord(Record{
		Level:       TRACE,
		MsgOrFormat: msg,
		Fields:      fields,
		Stack:       Enable,
		StackSize:   defStackSize,
	})
}

func (d DynamicLogger) IsEnabled(level Level) bool {
	return getDefLogger().IsEnabled(level)
}

func (d DynamicLogger) log(r Record) {
	r.CallerSkip++
	getDefLogger().log(r)
}

func (d DynamicLogger) buildFields(fields ...Field) []Field {
	return fields
}

func (d DynamicLogger) printRecord(r Record) {
	l := getDefLogger()
	if l.IsEnabled(r.Level) {
		r.Caller = d.Caller
		r.ShortFile = d.ShortFile
		r.CallerSkip = defCallerSkip + d.CallerSkip
		r.App = d.app
		l.output(r)
	}
}
