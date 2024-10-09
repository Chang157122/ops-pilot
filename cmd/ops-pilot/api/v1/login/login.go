package login

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/common/request"
	"opsPilot/internal/login"
	"opsPilot/internal/pkg/e"
)

// ValidateCodeLoginController 验证码登录接口
func ValidateCodeLoginController(c *gin.Context) {
	var (
		loginRequest request.LoginRequest
	)
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		e.Error(c, e.InvalidParams)
	}
	e.Success(c, login.ValidateCodeLogin(loginRequest))
}
