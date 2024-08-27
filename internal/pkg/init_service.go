package pkg

import (
	"opsPilot/internal/pkg/conf"
	"opsPilot/internal/pkg/log"
)

func Init() {
	// 配置初始化
	conf.InitConfig()

	// 日志初始化
	log.InitLogger()
}
