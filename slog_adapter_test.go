//go:build go1.21

package olog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/welllog/olog/encoder"
)

func TestSlogHandlerWithAttrs(t *testing.T) {
	fields := []Field{
		{
			Key:   "k1",
			Value: "v1",
		},
	}

	validate := func(fds []Field) func(vfds []Field, t *testing.T) {
		return func(vfds []Field, t *testing.T) {
			if len(fds) != len(vfds) {
				t.Fatalf("expected %d fields, got %d", len(vfds), len(fds))
			}
			for i, vfd := range vfds {
				if fds[i].Key != vfd.Key {
					t.Fatalf("expected key %s, got %s", vfd.Key, fds[i].Key)
				}
				if fds[i].Value != vfd.Value {
					t.Fatalf("expected value %s, got %s", vfd.Value, fds[i].Value)
				}
			}
		}
	}

	vfn1 := validate(fields)
	handler := NewSlogHandler(NewLogger(
		WithLoggerWriter(NewWriter(io.Discard)),
		WithLoggerEncodeFunc(func(record Record, buffer *encoder.Buffer) {
			vfn1(record.Fields, t)
		}),
	))
	logger := slog.New(handler)
	logger.With(slog.String("k1", "v1")).Info("yes")

	fields = []Field{
		{
			Key:   "p1.k2",
			Value: "v2",
		},
		{
			Key:   "p1.k3",
			Value: "v3",
		},
		{
			Key:   "p1.p2.k4",
			Value: "v4",
		},
		{
			Key:   "p1.k5",
			Value: "v5",
		},
	}
	vfn2 := validate(fields)
	handler = NewSlogHandler(NewLogger(
		WithLoggerWriter(NewWriter(io.Discard)),
		WithLoggerEncodeFunc(func(record Record, buffer *encoder.Buffer) {
			vfn2(record.Fields, t)
		}),
	))
	logger = slog.New(handler)
	logger.WithGroup("p1").With(
		slog.Group("p2", slog.String("k4", "v4")),
		slog.String("k5", "v5"),
	).
		Info("yes", slog.String("k2", "v2"), slog.String("k3", "v3"))
}

func TestSlogCaller(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := slog.New(
		NewSlogHandler(
			NewLogger(WithLoggerWriter(NewWriter(buf))),
		),
	)

	ctx := context.Background()
	logger.Log(ctx, slog.LevelInfo, "hello")
	logger.Error("hello")
	logger.ErrorContext(ctx, "hello")
	logger.Warn("hello")
	logger.WarnContext(ctx, "hello")
	logger.Info("hello")
	logger.InfoContext(ctx, "hello")
	logger.Debug("hello")
	logger.DebugContext(ctx, "hello")
	logger.With(slog.String("k1", "v1")).WithGroup("p1").With(slog.String("k2", "v2")).Info("yes")

	err := validCaller(buf, "olog/slog_adapter_test.go", 92)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkSlogHandler(b *testing.B) {
	logger1 := slog.New(
		slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			AddSource: true,
		}),
	)
	logger2 := slog.New(
		NewSlogHandler(NewLogger(
			WithLoggerWriter(NewWriter(io.Discard)),
		)),
	)
	logger3 := NewLogger(
		WithLoggerWriter(NewWriter(io.Discard)),
	)

	b.Run("slog", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger1.Info("hello")
		}
	})
	b.Run("slog_fields", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger1.Info("hello", slog.String("k1", "v1"), slog.Int("k2", 2))
		}
	})
	b.Run("slog_infof", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger1.Info(fmt.Sprintf("hello %s", "world"))
		}
	})

	b.Run("slog_h", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger2.Info("hello")
		}
	})
	b.Run("slog_h_fields", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger2.Info("hello", slog.String("k1", "v1"), slog.Int("k2", 2))
		}
	})
	b.Run("slog_h_infof", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger2.Info(fmt.Sprintf("hello %s", "world"))
		}
	})

	b.Run("olog", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger3.Info("hello")
		}
	})
	b.Run("olog_fields", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger3.Infow("hello", Field{Key: "k1", Value: "v1"}, Field{Key: "k2", Value: 2})
		}
	})
	b.Run("olog_infof", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger3.Infof("hello %s", "world")
		}
	})
}
