package logger

import (
	"strconv"
	"testing"
)

func TestFileLogger(t *testing.T) {
	config := make(map[string]string, 4)

	config["log_level"] = strconv.Itoa(LogLevelDebug)
	config["log_path"] =  "/Users/wujian/go/src/homework/logs"
	config["log_name"] =  "test"

	logger,_ := NewFileLogger(config)

	logger.Debug("user id [%d] is come from china", 123)
	logger.Info("user id [%d] is come from china", 123)
	logger.Warn("user id [%d] is come from china", 123)
	logger.Fatal("user id [%d] is come from china", 123)
}

func TestConsoleLogger(t *testing.T) {
	config := make(map[string]string, 4)
	config["log_level"] = strconv.Itoa(LogLevelDebug)
	config["log_path"] =  "/Users/wujian/go/src/homework/logs"
	config["log_name"] =  "test"

	logger,_ := NewFileLogger(config)

	logger.Debug("user id [%d] is come from china", 123)
}
