package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	Service *ServiceConfig `json:"service"`
}

var (
	configFile string
	C          *Configuration
)

func InitConfig() {
	if os.Getenv("CONFIG_FILE") == "" {
		configFile = "conf/application.yaml"
	}
	open, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(open, &C); err != nil {
		panic(err)
	}
}
