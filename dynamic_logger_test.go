package olog

import (
	"bytes"
	"context"
	"testing"
)

func TestDynamicLoggerCaller(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	SetLoggerOptions(
		WithLoggerCaller(true),
		WithLoggerWriter(NewWriter(buf)),
	)

	logger := DynamicLogger{}
	logger.Log(Record{MsgOrFormat: "hello"})
	logger.Error("hello")
	logger.Errorf("hello")
	logger.Errorw("hello")
	logger.Warn("hello")
	logger.Warnf("hello")
	logger.Warnw("hello")
	logger.Notice("hello")
	logger.Noticef("hello")
	logger.Noticew("hello")
	logger.Info("hello")
	logger.Infof("hello")
	logger.Infow("hello")
	logger.Debug("hello")
	logger.Debugf("hello")
	logger.Debugw("hello")
	logger.Trace("hello")
	logger.Tracef("hello")
	logger.Tracew("hello")
	WithContext(logger, context.Background()).Log(Record{MsgOrFormat: "hello"})
	logger.log(Record{MsgOrFormat: "hello", CallerSkip: -1})

	err := validCaller(buf, "olog/dynamic_logger_test.go", 17)
	if err != nil {
		t.Error(err)
	}
}
