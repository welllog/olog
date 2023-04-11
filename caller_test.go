package olog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLogCaller(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	SetWriter(NewWriter(buf))
	SetEncode(JSON)
	SetLevel(DEBUG)

	Log(DEBUG, WithPrintMsg("hello"))
	Error("hello")
	Errorf("hello %s", "world")
	Errorw("hello")
	Warn("hello")
	Warnf("hello %s", "world")
	Warnw("hello")
	Info("hello")
	Infof("hello %s", "world")
	Infow("hello")
	Debug("hello")
	Debugf("hello %s", "world")
	Debugw("hello")

	err := validCaller(buf, "olog/caller_test.go", 17)
	if err != nil {
		t.Error(err)
	}
}

func TestLoggerCaller(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(WithLoggerWriter(NewWriter(buf)))

	logger.Log(DEBUG, WithPrintMsg("hello"))
	logger.Error("hello")
	logger.Errorf("hello %s", "world")
	logger.Errorw("hello")
	logger.Warn("hello")
	logger.Warnf("hello %s", "world")
	logger.Warnw("hello")
	logger.Info("hello")
	logger.Infof("hello %s", "world")
	logger.Infow("hello")
	logger.Debug("hello")
	logger.Debugf("hello %s", "world")
	logger.Debugw("hello")

	err := validCaller(buf, "olog/caller_test.go", 41)
	if err != nil {
		t.Error(err)
	}
}

func TestContextLoggerCaller(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(WithLoggerWriter(NewWriter(buf)))
	logger = WithContext(logger, context.Background())
	logger = WithContext(logger, context.WithValue(context.Background(), "name", "bob"))

	logger.Log(DEBUG, WithPrintMsg("hello"))
	logger.Error("hello")
	logger.Errorf("hello %s", "world")
	logger.Errorw("hello")
	logger.Warn("hello")
	logger.Warnf("hello %s", "world")
	logger.Warnw("hello")
	logger.Info("hello")
	logger.Infof("hello %s", "world")
	logger.Infow("hello")
	logger.Debug("hello")
	logger.Debugf("hello %s", "world")
	logger.Debugw("hello")

	err := validCaller(buf, "olog/caller_test.go", 67)
	if err != nil {
		t.Error(err)
	}
}

func validCaller(buf *bytes.Buffer, file string, startLine int) error {
	for {
		line, _ := buf.ReadBytes('\n')
		if len(line) == 0 {
			break
		}
		m := make(map[string]any)
		if err := json.Unmarshal(line, &m); err != nil {
			return fmt.Errorf("unmarshal data %s err: %w", string(line), err)
		}
		want := fmt.Sprintf("%s:%d", file, startLine)
		if m[fieldCaller] != want {
			return fmt.Errorf("caller is %s, want %s", m[fieldCaller], want)
		}
		startLine++
	}
	return nil
}
