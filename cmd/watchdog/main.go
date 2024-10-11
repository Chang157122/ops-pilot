package main

import (
	"opsPilot/cmd/watchdog/api"
)

func main() {
	router := api.InitRouter()
	if err := router.Run(":7798"); err != nil {
		panic(err)
	}
}
