package logger

import "testing"

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger(LogLevelDebug, "/Users/wujian/go/src/homework/logs", "test")

	logger.Debug("user id [%d] is come from china", 123)
	logger.Info("user id [%d] is come from china", 123)
	logger.Warn("user id [%d] is come from china", 123)
	logger.Fatal("user id [%d] is come from china", 123)
}

func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(LogLevelDebug)

	logger.Debug("user id [%d] is come from china", 123)
}
