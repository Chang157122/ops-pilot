package conf

import "time"

// ServiceConfig 基础配置信息
type ServiceConfig struct {
	Mode          string        `yaml:"mode"`
	LogFilePath   string        `yaml:"log_file_path,logFilePath"`
	LogFileMaxAge time.Duration `yaml:"log_file_max_age,logFileMaxAge"`
}
