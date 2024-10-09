package api

import (
	"github.com/gin-gonic/gin"
	v1 "opsPilot/cmd/watchdog/api/v1"
	"opsPilot/internal/pkg/middleeware"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()
	gin.SetMode("release")
	engine.Use(middleeware.GinRecovery(true), middleeware.GinLogger())

	group := engine.Group("/watchdog")
	V1 := group.Group("/v1")
	{
		V1.GET("/metrics", v1.MetricsController)
	}
	return engine
}
