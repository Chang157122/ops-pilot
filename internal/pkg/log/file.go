package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

type FileLogger struct {
	// 日志级别 比该级别高的打印的时候会显示
	Level LogLevel `json:"level"`
	// 前缀 打印日志信息的前缀
	Prefix string `json:"prefix"`
	// 日志文件夹路径
	FilePath string `json:"file_path"`
	// 日志文件名
	FileNameList []string `json:"file_name"`
	//// 错误日志文件名
	//ErrFileName string `json:"err_file_name"`
	// 最大文件大小 分割日志文件的时候会用到
	MaxFileSize int64 `json:"max_file_size"`
	// 最大文件保留个数
	MaxFileCount int `json:"max_file_count"`
	// 分割文件最大时间间隔 按时间分割文件的时候会用到
	SplitFileAge int `json:"split_file_age"`
	// 最大保存保存时间长度
	MaxFileAge int `json:"max_file_age"`
	// 是否打印文件行号信息
	FileInfoLevel LogLevel `json:"file_info_level"`
	// 输出类型
	Encoder encoder `json:"encoder"`
	// 记录上次切割文件的时间
	LastSplitTime time.Time
	// 日志文件对象
	FileObj map[string]*os.File
	// 异步写日志 管道
	LogChan   chan *LogMsg
	CheckTime time.Duration
	// 加锁,用于日志分割
	lock sync.RWMutex
}

type LogMsg struct {
	Level        LogLevel  `json:"level"`
	Msg          string    `json:"msg"`
	LogTime      time.Time `json:"log_time"`
	FileLineInfo string    `json:"file_line_info"`
}

// NewFileLogger 构造函数
func NewFileLogger(level LogLevel, fp string, options ...LogoHandleFunc) *FileLogger {

	fl := &FileLogger{
		Level:         level,
		FilePath:      fp,
		FileNameList:  []string{"access", "info", "error", "warning", "debug"},
		FileObj:       map[string]*os.File{},
		LastSplitTime: time.Now(),
		FileInfoLevel: FatalLevel,
		MaxFileCount:  1,
		CheckTime:     time.Hour,
		LogChan:       make(chan *LogMsg, maxChanSize),
		lock:          sync.RWMutex{},
	}
	for _, option := range options {
		if err := option.Configure(fl); err != nil {
			panic(err)
		}
	}
	err := fl.initFile()
	if err != nil {
		panic(err)
	}
	return fl
}

func WithMaxFileSize(size int64) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		fl.MaxFileSize = size
		return nil
	})
}

func WithSplitAge(splitFileAge int) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		fl.SplitFileAge = splitFileAge
		return nil
	})
}

func WithMaxAge(maxAge int) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		if fl.MaxFileCount > 0 && maxAge > 0 {
			return errors.New("attempt to set MaxAge when RotationCount is also given")
		}
		fl.MaxFileAge = maxAge
		return nil
	})
}

func WithMaxCount(count int) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		if fl.MaxFileAge > 0 && count > 0 {
			return errors.New("attempt to set RotationCount when MaxAge is also given")
		}
		if count <= 0 {
			return errors.New("RotationCount can't <= 0")
		}
		fl.MaxFileCount = count
		return nil
	})
}

func WithEncoder(encode encoder) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		fl.Encoder = encode
		return nil
	})
}

func WithCheckTime(t time.Duration) LogoHandleFunc {
	return OptionFn(func(fl *FileLogger) error {
		if fl.MaxFileCount <= 0 && fl.MaxFileAge <= 0 {
			return errors.New("attempt to set CheckTime when MaxFileCount <= 0 AND MaxFileAge <= 0")
		}
		fl.CheckTime = t
		return nil
	})
}

func (fl *FileLogger) initFile() error {

	// 判断目录不存在,则创建目录
	isExistDir(fl.FilePath)

	for _, fileName := range fl.FileNameList {
		fullFileName := path.Join(fl.FilePath, fileName+".log")
		fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fl.FileObj[fileName] = fileObj
		go fl.writeLogBackground()
		if fl.MaxFileAge > 0 {
			go fl.CheckFileMaxAge()
		}
		if fl.MaxFileCount > 0 {
			go fl.CheckFileCount()
		}
	}
	return nil
}

func isExistDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在,创建目录
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func (fl *FileLogger) CheckFileSize(file *os.File) bool {
	if fl.MaxFileSize == 0 {
		return false
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return false
	}
	return fileInfo.Size() >= fl.MaxFileSize
}

func (fl *FileLogger) CheckSplitFileAge(now time.Time) bool {
	if fl.SplitFileAge == 0 {
		return false
	}
	return int(now.Sub(fl.LastSplitTime).Hours()) > fl.SplitFileAge
}

func (fl *FileLogger) CheckFileMaxAge() {
	// 加载时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	fl.RemoveMoreMaxAgeFile(loc)
	ticker := time.Tick(fl.CheckTime)
	for _ = range ticker {
		fl.RemoveMoreMaxAgeFile(loc)
	}
}

func (fl *FileLogger) RemoveMoreMaxAgeFile(loc *time.Location) {
	earlist_time := fl.LastSplitTime.Add(-(time.Duration(fl.MaxFileAge) * time.Hour))
	file_list, _ := ioutil.ReadDir(fl.FilePath)
	for _, file := range file_list {
		if file.Name() != getLogFileName(file.Name()) {
			file_time := strings.Split(file.Name(), ".")[0]
			timeObj, err := time.ParseInLocation("20060102", file_time, loc)
			if err != nil {
				fl.Errorf("parse time error:%v", err)
				continue
			}
			if timeObj.Before(earlist_time) {
				if err := os.Remove(path.Join(fl.FilePath, file.Name())); err != nil {
					fl.Errorf("remove file failed:%v", err)
					continue
				}
				fl.Infof("%s 已过期 删除成功\n", file.Name())
			}
		}
	}
}

func (fl *FileLogger) CheckFileCount() {
	// 加载时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fl.Error(err.Error())
		return
	}
	fl.lock.Lock()
	defer fl.lock.Unlock()
	fl.RemoveMoreMaxCountFile(loc)
	//ticker := time.Tick(fl.CheckTime)
	//for _ = range ticker {
	//	fl.RemoveMoreMaxCountFile(loc)
	//}
}

func (fl *FileLogger) RemoveMoreMaxCountFile(loc *time.Location) {
	if fl.MaxFileCount == 0 {
		return
	}

	file_list, _ := ioutil.ReadDir(fl.FilePath)
	var logListMap = make(map[string][]os.FileInfo)
	for _, file := range file_list {
		if strings.HasSuffix(file.Name(), ".log") && file.Name() != getLogFileName(file.Name()) {
			fileName := strings.Split(file.Name(), ".")[0]
			fileTime := strings.Split(strings.Split(file.Name(), ".")[1], "-")[0]
			_, err := time.ParseInLocation("20060102", fileTime, loc)
			if err != nil {
				fl.Errorf("parse file %s as time error:%v", file.Name(), err)
				continue
			}
			logListMap[fileName] = append(logListMap[fileName], file)
		}
	}
	for _, v := range logListMap {
		if len(v) > fl.MaxFileCount {
			sort.SliceStable(v, func(i, j int) bool {
				return v[i].ModTime().After(v[j].ModTime())
			})
			for _, file := range v[fl.MaxFileCount:] {
				fmt.Println("fileName: ", file.Name())
				if err := os.Remove(path.Join(fl.FilePath, file.Name())); err != nil {
					fl.Errorf("remove file failed:%v", err)
					continue
				}
				fl.Infof("日志文件超出最大保留日志文件个数 %s 删除成功\n", file.Name())
			}
		}
	}
}

func (fl *FileLogger) SetPrefix(prefix string) {
	fl.Prefix = prefix
}

func (fl *FileLogger) SetTimeFormat(format string) {
	logTimeFormat = time.Now().Format(format)
}

func (fl *FileLogger) SetFileInfo(level LogLevel) {
	fl.FileInfoLevel = level
}

func (fl *FileLogger) SetEncoder(encode encoder) {
	fl.Encoder = encode
}

func (fl *FileLogger) Log(level LogLevel, format string, args ...interface{}) {
	if fl.Level > level {
		return
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	logTime := time.Now()
	info, err := getInfo(3)
	var log_msg *LogMsg
	if fl.FileInfoLevel <= level {
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		log_msg = &LogMsg{
			Level:        level,
			Msg:          msg,
			LogTime:      logTime,
			FileLineInfo: info,
		}
	} else {
		log_msg = &LogMsg{
			Level:        level,
			Msg:          msg,
			LogTime:      logTime,
			FileLineInfo: info,
		}
	}
	select {
	case fl.LogChan <- log_msg:
	default:
	}
}

func (fl *FileLogger) writeLogBackground() {
	for {
		select {
		case log_chan := <-fl.LogChan:
			if fl.CheckFileSize(fl.getLogFileObj(log_chan.Level)) || fl.CheckSplitFileAge(log_chan.LogTime) {
				newFileObj, err := fl.splitFile(fl.getLogFileObj(log_chan.Level))
				if err != nil {
					return
				}
				fl.LastSplitTime = log_chan.LogTime

				fl.FileObj[ParseLogLevel(log_chan.Level)] = newFileObj
			}
			var msg []byte
			if fl.Encoder == JsonEncoder {
				msg = fl.JsonEncode(log_chan)
			} else {
				msg = fl.TextEncode(log_chan)
			}
			io.Writer.Write(fl.FileObj[ParseLogLevel(log_chan.Level)], msg)
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (fl *FileLogger) TextEncode(logMsg *LogMsg) []byte {
	buf := bytes.Buffer{}
	if fl.Prefix != "" {
		buf.WriteString(fl.Prefix)
		buf.WriteString(" ")
	}
	buf.WriteString(logMsg.LogTime.Format(logTimeFormat))
	buf.WriteString(" ")
	if logMsg.FileLineInfo != "" {
		buf.WriteString(logMsg.FileLineInfo)
	}
	buf.WriteString(" ▶ [")
	buf.WriteString(logMsg.Level.String())
	buf.WriteString("] ")
	buf.WriteString(logMsg.Msg)
	buf.WriteString("\n")
	return buf.Bytes()
}

func (fl *FileLogger) JsonEncode(logMsg *LogMsg) []byte {
	buf := bytes.Buffer{}
	buf.WriteString(`{`)
	if fl.Prefix != "" {
		buf.WriteString(`"prefix": "`)
		buf.WriteString(fl.Prefix)
		buf.WriteString(`",`)
	}
	buf.WriteString(`"time": "`)
	buf.WriteString(logMsg.LogTime.Format(logTimeFormat))
	if logMsg.FileLineInfo != "" {
		buf.WriteString(`","fileInfo": "`)
		buf.WriteString(logMsg.FileLineInfo)
	}
	buf.WriteString(`","level": "`)
	buf.WriteString(logMsg.Level.String())
	buf.WriteString(`","msg": "`)
	buf.WriteString(logMsg.Msg)
	buf.WriteString("\"}\n")
	return buf.Bytes()
}

func (fl *FileLogger) splitFile(file *os.File) (*os.File, error) {
	var newLogName string

	nowStr := time.Now().Format("20060102")
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileName := fileInfo.Name()
	logName := path.Join(fl.FilePath, fileName)
	newFileName := strings.Replace(fileName, strings.Split(fileName, ".")[1], nowStr+".log", 1)
	newLogName = fl.isLogFileExist(newFileName, nowStr, 1)
	// 1. 关闭当前文件
	_ = file.Close()
	// 2. 备份一个 rename
	_ = os.Rename(logName, newLogName)
	// 3. 打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	// 4. 将打开的文件赋值给 fl.FileObj
	return fileObj, nil
}

// 切片文件是否存在
func (fl *FileLogger) isLogFileExist(fileName, nowStr string, i int) string {
	_, err := os.Stat(path.Join(fl.FilePath, fileName))
	// 文件不存在
	if err != nil {
		return path.Join(fl.FilePath, fileName)
	}
	logSplit := fmt.Sprintf("%s-%d", nowStr, i)

	return fl.isLogFileExist(strings.Replace(fileName, strings.Split(fileName, ".")[1], logSplit, 1), nowStr, i+1)

}

func (fl *FileLogger) Debug(msg string) {
	fl.Log(DebugLevel, msg)
}

func (fl *FileLogger) Trace(msg string) {
	fl.Log(TraceLevel, msg)
}

func (fl *FileLogger) Info(msg string) {
	fl.Log(InfoLevel, msg)
}

func (fl *FileLogger) Warning(msg string) {
	fl.Log(WarningLevel, msg)
}

func (fl *FileLogger) Error(msg string) {
	fl.Log(ErrorLevel, msg)
}

func (fl *FileLogger) Fatal(msg string) {
	fl.Log(FatalLevel, msg)
	os.Exit(1)
}

func (fl *FileLogger) Access(format string, args ...interface{}) {
	fl.Log(AccessLevel, format, args...)
}

func (fl *FileLogger) Debugf(format string, args ...interface{}) {
	fl.Log(DebugLevel, format, args...)
}

func (fl *FileLogger) Tracef(format string, args ...interface{}) {
	fl.Log(TraceLevel, format, args...)
}

func (fl *FileLogger) Infof(format string, args ...interface{}) {
	fl.Log(InfoLevel, format, args...)
}

func (fl *FileLogger) Warningf(format string, args ...interface{}) {
	fl.Log(WarningLevel, format, args...)
}

func (fl *FileLogger) Errorf(format string, args ...interface{}) {
	fl.Log(ErrorLevel, format, args...)
}
func (fl *FileLogger) Fatalf(format string, args ...interface{}) {
	fl.Log(FatalLevel, format, args...)
	time.Sleep(time.Second)
	os.Exit(1)
}
