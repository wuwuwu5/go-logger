package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

// 获取方法执行的位置
func GetLineInfo() (string, string, int) {
	// 获取程序执行的地方
	pc, file, line, ok := runtime.Caller(4) // 栈的深度
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		return file, funcName, line
	}

	return "", "", 0
}

func WriteLog(level int, format string, args ...interface{}) *LogData {

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

	warnAndFatal := false

	if level == LogLevelWarn || level == LogLevelError || level == LogLevelFail {
		warnAndFatal = true
	}

	return &LogData{
		Message:      msg,
		TimeStr:      timeStr,
		LevelStr:     levelText,
		FileName:     fileName,
		FuncName:     funcName,
		LineNo:       line,
		WarnAndFatal: warnAndFatal,
	}
}
