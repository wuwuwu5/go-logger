package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level        int
	path         string
	name         string
	file         *os.File
	warnFile     *os.File
	LogSplitType int
	LogSplitSize int64
	LogDataChan  chan *LogData
	LogSplitTime int
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

	logSplitTypeStr, ok := config["log_split_type"]
	var logSplitSize int64
	var logSplitType int

	if !ok {
		logSplitTypeStr = "date"
		logSplitType = LogSplitTypeDate
	} else {
		if logSplitTypeStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]

			if !ok {
				logSplitSizeStr = strconv.Itoa(100 * 1024 * 1024)
			}

			logSplitSize, _ = strconv.ParseInt(logSplitSizeStr, 10, 64)
			logSplitType = LogSplitTypeSize
		} else {
			logSplitType = LogSplitTypeDate
		}
	}

	// 管道大小
	logChanSize, ok := config["log_chan_size"]

	if !ok {
		logChanSize = "1024"
	}

	size, _ := strconv.Atoi(logChanSize)

	logger := &FileLogger{
		level:        getLevel(logLevel),
		path:         logPath,
		name:         logName,
		LogDataChan:  make(chan *LogData, size),
		LogSplitType: logSplitType,
		LogSplitSize: logSplitSize,
		LogSplitTime: time.Now().Hour(),
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

			this.checkSplitFile(v)
		}
	}()
}

// 根据时间分割日志
func (this *FileLogger) splitLogFileDate(logData *LogData) {
	now := time.Now()
	hour := now.Hour()

	if hour == this.LogSplitTime {
		return
	}

	var backupFileName string
	var olFileName string

	if logData.WarnAndFatal {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d", this.path, this.name, now.Year(), now.Month(), now.Day(), this.LogSplitTime)
		olFileName = fmt.Sprintf("%s/%s.log.wf", this.path, this.name)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d%02d", this.path, this.name, now.Year(), now.Month(), now.Day(), this.LogSplitTime)
		olFileName = fmt.Sprintf("%s/%s.log", this.path, this.name)
	}

	this.LogSplitTime = hour

	file := this.file

	if logData.WarnAndFatal {
		file = this.warnFile
	}

	file.Close()
	os.Rename(olFileName, backupFileName)

	file, err := os.OpenFile(olFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		return
	}

	if logData.WarnAndFatal {
		this.warnFile = file
	} else {
		this.file = file
	}
}

// 根据日志大小进行分割
func (this *FileLogger) splitLogFileSize(logData *LogData) {
	file := this.file
	if logData.WarnAndFatal {
		file = this.warnFile
	}

	info, err := file.Stat()

	if err != nil {
		return
	}

	fileSize := info.Size()

	if fileSize <= this.LogSplitSize {
		return
	}

	now := time.Now()

	var backupFileName string
	var olFileName string

	if logData.WarnAndFatal {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d%02d%02d", this.path, this.name, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		olFileName = fmt.Sprintf("%s/%s.log.wf", this.path, this.name)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d%02d%02d", this.path, this.name, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		olFileName = fmt.Sprintf("%s/%s.log", this.path, this.name)
	}

	file.Close()
	os.Rename(olFileName, backupFileName)

	file, err = os.OpenFile(olFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		return
	}

	if logData.WarnAndFatal {
		this.warnFile = file
	} else {
		this.file = file
	}
}

// 日志分割逻辑
func (this *FileLogger) checkSplitFile(logData *LogData) {
	if this.LogSplitType == LogSplitTypeDate {
		this.splitLogFileDate(logData)
	} else {
		this.splitLogFileSize(logData)
	}
}

// 设置等级
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
