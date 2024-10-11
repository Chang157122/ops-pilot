package main

import (
	"opsPilot/cmd/ops-monitor/api"
)

func main() {
	router := api.InitRouter()
	if err := router.Run(":7798"); err != nil {
		panic(err)
	}
}
