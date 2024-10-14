package settings

import "C"
import (
	"gopkg.in/yaml.v3"
	"opsPilot/internal/pkg/log"
	"os"
)

func InitConfig() {
	var CONFIG_FILE string
	if CONFIG_FILE = os.Getenv("CONFIG_FILE"); CONFIG_FILE == "" || len(CONFIG_FILE) == 0 {
		CONFIG_FILE = `application.yml`
	}
	file, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(file, &Cfg); err != nil {
		panic(err)
	}
	// 日志初始化
	func() {
		logConfig := Cfg.Service.LogConfig
		if err := log.InitLogger(logConfig.FileName, logConfig.Level, logConfig.MaxSize, logConfig.MaxBackups, logConfig.MaxAge); err != nil {
			panic(err)
		}
	}()
}

type Settings struct {
	Service ServiceConf `yaml:"service"`
}

type ServiceConf struct {
	Mode      string    `yaml:"mode"`
	LogConfig LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"file_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}
