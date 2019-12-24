package logger

import (
	"fmt"
	"os"
)

type FileLogger struct {
	level    int
	path     string
	name     string
	file     *os.File
	warnFile *os.File
}

func NewFileLogger(config map[string]string) (LogInterface, error) {

	log_level, ok := config["log_level"]

	if !ok {
		return nil, fmt.Errorf("not found log_level")
	}

	log_path, ok := config["log_path"]

	if !ok {
		return nil, fmt.Errorf("not found log_path")
	}

	log_name, ok := config["log_name"]

	if !ok {
		return nil, fmt.Errorf("not found log_name")
	}

	logger := &FileLogger{
		level: getLevel(log_level),
		path:  log_path,
		name:  log_name,
	}

	logger.Init()

	return logger, nil
}

func (this *FileLogger) Init() {
	fileName := fmt.Sprintf("%s/%s.log", this.path, this.name)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		panic(fmt.Sprintf("open %s faile, err: %v", fileName, err))
	}

	this.file = file

	// error和fail级别的日志单独写
	fileName = fmt.Sprintf("%s/%s.log.wf", this.path, this.name)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		panic(fmt.Sprintf("open %s faile, err: %v", fileName, err))
	}

	this.warnFile = file

}

func (this *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFail {
		this.level = LogLevelDebug
	} else {
		this.level = level
	}
}

func (this *FileLogger) Debug(format string, args ...interface{}) {
	if this.level > LogLevelDebug {
		return
	}
	WriteLog(this.file, LogLevelDebug, format, args...)
}

func (this *FileLogger) Trace(format string, args ...interface{}) {
	if this.level > LogLevelTrace {
		return
	}
	WriteLog(this.file, LogLevelTrace, format, args...)
}

func (this *FileLogger) Info(format string, args ...interface{}) {
	if this.level > LogLevelInfo {
		return
	}
	WriteLog(this.file, LogLevelInfo, format, args...)
}

func (this *FileLogger) Warn(format string, args ...interface{}) {
	if this.level > LogLevelWarn {
		return
	}
	WriteLog(this.file, LogLevelWarn, format, args...)
}

func (this *FileLogger) Error(format string, args ...interface{}) {
	if this.level > LogLevelError {
		return
	}
	WriteLog(this.warnFile, LogLevelError, format, args...)
}

func (this *FileLogger) Fatal(format string, args ...interface{}) {
	if this.level > LogLevelFail {
		return
	}
	WriteLog(this.warnFile, LogLevelFail, format, args...)
}

func (this *FileLogger) Close() {
	this.file.Close()
	this.warnFile.Close()
}
