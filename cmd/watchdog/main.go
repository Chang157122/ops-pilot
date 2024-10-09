package watchdog

import (
	"opsPilot/cmd/watchdog/api"
)

func APP() {
	router := api.InitRouter()
	if err := router.Run(":7798"); err != nil {
		panic(err)
	}
}
