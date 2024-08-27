package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type LogLevel uint8
type encoder uint8

type Logger interface {
	Debug(msg string)
	Trace(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Fatal(msg string)
	Access(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	SetPrefix(prefix string)
	SetTimeFormat(format string)
	SetFileInfo(level LogLevel)
	SetEncoder(encode encoder)
}

type MaxFileCount struct {
	FileCount    int
	ErrFileCount int
}

//type MaxFileSize struct {
//	Size int64
//}

//type FileAge struct {
//	SplitFileAge int
//	MaxFileAge int
//}

const (
	UnKnow LogLevel = iota
	DebugLevel
	TraceLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	AccessLevel
	FatalLevel
)

const (
	TextEncoder encoder = iota
	JsonEncoder
)

var (
	levels = map[LogLevel]string{
		DebugLevel:   "Debug",
		AccessLevel:  "Access",
		TraceLevel:   "Trace",
		InfoLevel:    "Info",
		WarningLevel: "Warning",
		ErrorLevel:   "Error",
		FatalLevel:   "Fatal",
	}
	// 日志时间格式字符串
	logTimeFormat = "2006/01/02  15:04:05.000"
	// 是否打印文件行号信息 吗，偶人为true
	maxChanSize = 50000
)

func (ll LogLevel) String() string {
	return levels[ll]
}

func ParseLogLevel(level LogLevel) string {
	switch level {
	case DebugLevel:
		return "debug"
	case AccessLevel:
		return "access"
	case InfoLevel:
		return "info"
	case WarningLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}

func getInfo(n int) (string, error) {
	pc, fileName, lineNo, ok := runtime.Caller(n)
	if !ok {
		return "", errors.New("runtime.Caller() failed\n")
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	//return fmt.Sprintf("%s:%s:%d", fileName, funcName, lineNo), nil
	return fmt.Sprintf("%s:%d", fileName, lineNo), nil
}

func New() *ConsoleLogger {
	return NewConsoleLogger(DebugLevel)
}

func getLogFileName(fileName string) string {
	switch fileName {
	case "access.log":
		return "access.log"
	case "debug.log":
		return "debug.log"
	case "info.log":
		return "info.log"
	case "sql.log":
		return "sql.log"
	case "warning.log":
		return "warning.log"
	case "error.log":
		return "error.log"
	default:
		return "unknown.log"
	}
}

func (fl *FileLogger) getLogFileObj(level LogLevel) *os.File {
	switch level {
	case DebugLevel:
		return fl.FileObj["debug"]
	case InfoLevel:
		return fl.FileObj["info"]
	case WarningLevel:
		return fl.FileObj["warning"]
	case ErrorLevel:
		return fl.FileObj["error"]
	case AccessLevel:
		return fl.FileObj["access"]
	default:
		return fl.FileObj["info"]
	}
}
