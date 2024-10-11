package login

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/common/request"
	login2 "opsPilot/internal/ops-pilot/login"
	"opsPilot/internal/pkg/e"
)

var login = login2.NewLoginInterFace()

// GenCaptchaController 生成验证码
func GenCaptchaController(c *gin.Context) {
	e.Success(c, login.CaptchaGenerate())
}

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

// RegisteredUserController 注册用户
func RegisteredUserController(c *gin.Context) {
	var (
		req request.RegisterUserRequest
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(c, e.InvalidParams)
	}
	e.Success(c, login.RegisterUser(req))
}
