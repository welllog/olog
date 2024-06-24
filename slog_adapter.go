//go:build go1.21

package olog

import (
	"context"
	"log/slog"
)

var slogLevel2Level = map[slog.Level]Level{
	slog.LevelDebug: DEBUG,
	slog.LevelInfo:  INFO,
	slog.LevelWarn:  WARN,
	slog.LevelError: ERROR,
}

type SlogHandler struct {
	prefix    string
	logger    Logger
	ctxHandle CtxHandle
}

func NewSlogHandler(logger Logger, handles ...CtxHandle) SlogHandler {
	var handle CtxHandle
	if len(handles) > 0 {
		handle = handles[0]
	} else {
		handle = getDefCtxHandle()
	}

	return SlogHandler{
		logger:    logger,
		ctxHandle: handle,
	}
}

func (s SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return s.logger.IsEnabled(slogLevel2Level[level])
}

func (s SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := s.ctxHandle(ctx)

	r := Record{
		Level:       slogLevel2Level[record.Level],
		CallerSkip:  2,
		MsgOrFormat: record.Message,
		Fields:      s.logger.buildFields(fields...),
		Time:        record.Time,
	}

	attrLen := record.NumAttrs()
	if attrLen > 0 {
		fields = make([]Field, 0, len(r.Fields)+attrLen)
		record.Attrs(func(attr slog.Attr) bool {
			fields = addAttrsToFields(fields, s.prefix, attr)
			return true
		})

		fields = append(fields, r.Fields...)
		r.Fields = fields
	}

	s.logger.log(r)

	return nil
}

func (s SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]Field, 0, len(attrs))
	fields = addAttrsToFields(fields, s.prefix, attrs...)

	return SlogHandler{
		prefix: s.prefix,
		logger: &ctxLogger{
			Logger: s.logger,
			fields: fields,
		},
		ctxHandle: s.ctxHandle,
	}
}

func (s SlogHandler) WithGroup(name string) slog.Handler {
	return SlogHandler{
		prefix:    attrKeyJoin(s.prefix, name),
		logger:    s.logger,
		ctxHandle: s.ctxHandle,
	}
}

func addAttrsToFields(fields []Field, prefix string, attrs ...slog.Attr) []Field {
	for _, attr := range attrs {
		key := attrKeyJoin(prefix, attr.Key)
		if attr.Value.Kind() == slog.KindGroup {
			fields = addAttrsToFields(fields, key, attr.Value.Group()...)
		} else {
			fields = append(fields, Field{
				Key:   key,
				Value: attr.Value.Any(),
			})
		}
	}

	return fields
}

func attrKeyJoin(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}
