package main

import "opsPilot/cmd/ops-pilot/api"

func main() {
	APP()
}

func APP() {
	router := api.InitRouter()
	if err := router.Run(":7980"); err != nil {
		panic(err)
	}
}
