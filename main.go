package main

import (
	"opsPilot/internal/pkg/log"
	"time"
)

func main() {
	logger := log.NewFileLogger(log.DebugLevel, "log",
		log.WithSplitAge(24*3600*7),
		log.WithCheckTime(time.Hour),
	)
	logger.SetEncoder(log.TextEncoder)
	for i := 0; i < 1000; i++ {
		logger.Info("debug 测试信息")
		logger.Errorf("debug 测试信息")

	}
	time.Sleep(2 * time.Second)
}
