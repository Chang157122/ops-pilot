package log

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"time"
)

var (
	Logger *zap.SugaredLogger
)

func init() {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	var core zapcore.Core

	// 开发模式，日志输出到终端
	//consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	// NewTee创建一个核心，将日志条目复制到两个或多个底层核心中。
	core = zapcore.NewTee(
		//zapcore.NewCore(encoder, zapcore.AddSync(appHook), level),
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
	)

	// 创建 logger 对象
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	Logger = log.Sugar()

	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(log)
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
