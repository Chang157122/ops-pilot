package api

import (
	"github.com/gin-gonic/gin"
	"opsPilot/cmd/ops-pilot/api/v1/login"
	"opsPilot/internal/pkg/middleeware"
)

func InitRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(middleeware.GinRecovery(true), middleeware.GinLogger())

	group := engine.Group("/ops-pilot")
	V1 := group.Group("/v1")
	{
		V1.GET("/captcha", login.GenCaptchaController)
		V1.POST("/validateCodeLogin", login.ValidateCodeLoginController)
		V1.POST("/registeredUser", login.RegisteredUserController)
	}

	return engine
}
