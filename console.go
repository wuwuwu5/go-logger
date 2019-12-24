package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config map[string]string) (*ConsoleLogger, error) {

	log_level, ok := config["log_level"]

	if !ok {
		return nil, fmt.Errorf("not found log_level")
	}

	return &ConsoleLogger{level: getLevel(log_level)}, nil
}

func (this *ConsoleLogger) Init() {

}

func (this *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFail {
		this.level = LogLevelDebug
	} else {
		this.level = level
	}
}

func (this *ConsoleLogger) writeLog(file *os.File, level int, format string, args ...interface{}) {
	if this.level > level {
		return
	}

	now := time.Now()

	// 获取时间
	timeStr := now.Format("2006-01-02 15:04:05.999")

	// 获取日志等级
	levelText := getLevelText(level)

	// 获取执行位置文件名称  方法名称 行数
	fileName, funcName, line := GetLineInfo()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)

	// 格式用户输入的内容
	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(file, "%s %s %s:%s:%d %s \n", timeStr, levelText, fileName, funcName, line, msg)
}

func (this *ConsoleLogger) Debug(format string, args ...interface{}) {
	if this.level > LogLevelDebug {
		return
	}

	logData := WriteLog(LogLevelDebug, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) Trace(format string, args ...interface{}) {
	if this.level > LogLevelTrace {
		return
	}

	logData := WriteLog(LogLevelTrace, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) Info(format string, args ...interface{}) {
	if this.level > LogLevelInfo {
		return
	}

	logData := WriteLog(LogLevelInfo, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) Warn(format string, args ...interface{}) {
	if this.level > LogLevelWarn {
		return
	}

	logData := WriteLog(LogLevelWarn, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) Error(format string, args ...interface{}) {
	if this.level > LogLevelError {
		return
	}

	logData := WriteLog(LogLevelError, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if this.level > LogLevelFail {
		return
	}

	logData := WriteLog(LogLevelFail, format, args...)

	this.console(logData)
}

func (this *ConsoleLogger) console(logData *LogData) {
	fmt.Fprintf(os.Stdout, "%s %s %s:%d %s %s\n", logData.TimeStr, logData.LevelStr, logData.FileName, logData.LineNo, logData.FuncName, logData.Message)
}

func (this *ConsoleLogger) Close() {

}
