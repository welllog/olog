package olog

type DynamicLogger struct{}

func (d DynamicLogger) Log(r Record) {
	l := getDefLogger()
	if l.IsEnabled(r.Level) {
		l.log(r)
	}
}

func (d DynamicLogger) Fatal(args ...any) {
	getDefLogger().fatal(args...)
}

func (d DynamicLogger) Fatalf(format string, args ...any) {
	getDefLogger().fatalf(format, args...)
}

func (d DynamicLogger) Fatalw(msg string, fields ...Field) {
	getDefLogger().fatalw(msg, fields...)
}

func (d DynamicLogger) Error(args ...any) {
	getDefLogger().error(args...)
}

func (d DynamicLogger) Errorf(format string, args ...any) {
	getDefLogger().errorf(format, args...)
}

func (d DynamicLogger) Errorw(msg string, fields ...Field) {
	getDefLogger().errorw(msg, fields...)
}

func (d DynamicLogger) Warn(args ...any) {
	getDefLogger().warn(args...)
}

func (d DynamicLogger) Warnf(format string, args ...any) {
	getDefLogger().warnf(format, args...)
}

func (d DynamicLogger) Warnw(msg string, fields ...Field) {
	getDefLogger().warnw(msg, fields...)
}

func (d DynamicLogger) Notice(args ...any) {
	getDefLogger().notice(args...)
}

func (d DynamicLogger) Noticef(format string, args ...any) {
	getDefLogger().noticef(format, args...)
}

func (d DynamicLogger) Noticew(msg string, fields ...Field) {
	getDefLogger().noticew(msg, fields...)
}

func (d DynamicLogger) Info(args ...any) {
	getDefLogger().info(args...)
}

func (d DynamicLogger) Infof(format string, args ...any) {
	getDefLogger().infof(format, args...)
}

func (d DynamicLogger) Infow(msg string, fields ...Field) {
	getDefLogger().infow(msg, fields...)
}

func (d DynamicLogger) Debug(args ...any) {
	getDefLogger().debug(args...)
}

func (d DynamicLogger) Debugf(format string, args ...any) {
	getDefLogger().debugf(format, args...)
}

func (d DynamicLogger) Debugw(msg string, fields ...Field) {
	getDefLogger().debugw(msg, fields...)
}

func (d DynamicLogger) Trace(args ...any) {
	getDefLogger().trace(args...)
}

func (d DynamicLogger) Tracef(format string, args ...any) {
	getDefLogger().tracef(format, args...)
}

func (d DynamicLogger) Tracew(msg string, fields ...Field) {
	getDefLogger().tracew(msg, fields...)
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
