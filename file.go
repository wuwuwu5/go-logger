package logger

import (
	"fmt"
	"os"
	"strconv"
)

type FileLogger struct {
	level       int
	path        string
	name        string
	file        *os.File
	warnFile    *os.File
	LogDataChan chan *LogData
}

// 储存日志信息结构体
type LogData struct {
	Message      string
	TimeStr      string
	LevelStr     string
	FileName     string
	FuncName     string
	LineNo       int
	WarnAndFatal bool
}

func NewFileLogger(config map[string]string) (LogInterface, error) {

	logLevel, ok := config["log_level"]

	if !ok {
		return nil, fmt.Errorf("not found log_level")
	}

	logPath, ok := config["log_path"]

	if !ok {
		return nil, fmt.Errorf("not found log_path")
	}

	logName, ok := config["log_name"]

	if !ok {
		return nil, fmt.Errorf("not found log_name")
	}

	logChanSize, ok := config["log_chan_size"]

	if !ok {
		logChanSize = "1024"
	}

	size, _ := strconv.Atoi(logChanSize)

	logger := &FileLogger{
		level:       getLevel(logLevel),
		path:        logPath,
		name:        logName,
		LogDataChan: make(chan *LogData, size),
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

	go func() {
		for v := range this.LogDataChan {
			file := this.file
			if v.WarnAndFatal {
				file = this.warnFile
			}

			fmt.Fprintf(file, "%s %s %s:%d %s %s\n", v.TimeStr, v.LevelStr, v.FileName, v.LineNo, v.FuncName, v.Message)
		}
	}()
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

	logData := WriteLog(LogLevelDebug, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Trace(format string, args ...interface{}) {
	if this.level > LogLevelTrace {
		return
	}

	logData := WriteLog(LogLevelTrace, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Info(format string, args ...interface{}) {
	if this.level > LogLevelInfo {
		return
	}
	logData := WriteLog(LogLevelInfo, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Warn(format string, args ...interface{}) {
	if this.level > LogLevelWarn {
		return
	}
	logData := WriteLog(LogLevelWarn, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Error(format string, args ...interface{}) {
	if this.level > LogLevelError {
		return
	}

	logData := WriteLog(LogLevelError, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Fatal(format string, args ...interface{}) {
	if this.level > LogLevelFail {
		return
	}

	logData := WriteLog(LogLevelFail, format, args...)

	this.dispath(logData)
}

func (this *FileLogger) Close() {
	this.file.Close()
	this.warnFile.Close()
}

func (this *FileLogger) dispath(logData *LogData) {
	// 用于判断队列是否满了
	select {
	case this.LogDataChan <- logData:
	default:
	}
}
