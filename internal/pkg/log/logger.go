package log

import "C"
import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.SugaredLogger

//func init() {
//	var cfg settings.Settings
//	if cfg.Service.Mode == "debug" {
//		InitConsoleLogger()
//	} else {
//		if err := InitLogger(cfg.Service.LogConfig.FileName,
//			cfg.Service.LogConfig.Level,
//			cfg.Service.LogConfig.MaxSize,
//			cfg.Service.LogConfig.MaxBackups,
//			cfg.Service.LogConfig.MaxAge); err != nil {
//			panic(err)
//		}
//	}
//}

func InitLogger(filename, level string, maxSize, maxBackup, maxAge int) (err error) {
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge)
	encoder := getEncoder()
	l := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	core := zapcore.NewCore(encoder, writeSyncer, l)

	log := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(log) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	Logger = log.Sugar()
	return
}

func init() {
	encoder := getEncoder()
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), debugLevel)
	log := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(log)
	Logger = log.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
