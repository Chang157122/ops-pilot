package api

import (
	"github.com/gin-gonic/gin"
	v1 "opsPilot/cmd/ops-monitor/api/v1"
	"opsPilot/internal/pkg/middleeware"
)

func InitRouter() *gin.Engine {
	engine := gin.New()
	gin.SetMode("debug")
	engine.Use(middleeware.GinRecovery(true), middleeware.GinLogger())

	group := engine.Group("/ops-monitor")
	V1 := group.Group("/v1")
	{
		V1.GET("/metrics", v1.MetricsController)
		V1.POST("/checkThirdPort", v1.CheckThirdPortController)
	}
	return engine
}
