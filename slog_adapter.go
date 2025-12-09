//go:build go1.21

package olog

import (
	"context"
	"log/slog"
	"sync"
)

var fieldPool = sync.Pool{
	New: func() any {
		var f []Field
		return &f
	},
}

func toLevel(l slog.Level) Level {
	switch l {
	case slog.LevelDebug:
		return DEBUG
	case slog.LevelInfo:
		return INFO
	case slog.LevelWarn:
		return WARN
	case slog.LevelError:
		return ERROR
	default:
		return INFO
	}
}

type SlogHandler struct {
	prefix    string
	logger    Logger
	ctxHandle CtxHandle
}

func NewSlogHandler(logger Logger, handles ...CtxHandle) *SlogHandler {
	var handle CtxHandle
	if len(handles) > 0 {
		handle = handles[0]
	} else {
		handle = getDefCtxHandle()
	}

	return &SlogHandler{
		logger:    logger,
		ctxHandle: handle,
	}
}

func (s *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return s.logger.IsEnabled(toLevel(level))
}

func (s *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := s.ctxHandle(ctx)

	r := Record{
		Level:       toLevel(record.Level),
		CallerSkip:  3,
		MsgOrFormat: record.Message,
		Fields:      fields,
		Time:        record.Time,
	}

	attrLen := record.NumAttrs()
	if attrLen > 0 {
		fp := fieldPool.Get().(*[]Field)
		flds := (*fp)[:0]
		if cap(flds) < attrLen+len(fields) {
			flds = make([]Field, 0, attrLen+len(fields))
		}

		record.Attrs(func(attr slog.Attr) bool {
			flds = addAttrsToFields(flds, s.prefix, attr)
			return true
		})

		flds = append(flds, r.Fields...)
		r.Fields = flds

		s.logger.Log(r)

		// free fields
		// for i := range flds {
		// 	flds[i].Key = ""
		// 	flds[i].Value = nil
		// }
		clear(flds)
		*fp = flds
		fieldPool.Put(fp)
		return nil
	}

	s.logger.Log(r)

	return nil
}

func (s *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]Field, 0, len(attrs))
	fields = addAttrsToFields(fields, s.prefix, attrs...)

	return &SlogHandler{
		prefix:    s.prefix,
		logger:    WithFields(s.logger, fields...),
		ctxHandle: s.ctxHandle,
	}
}

func (s *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{
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
