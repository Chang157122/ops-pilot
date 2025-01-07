package login

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/common/request"
	login2 "opsPilot/internal/ops-pilot/login"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/pkg/jwt"
	"opsPilot/internal/pkg/log"
	"opsPilot/internal/pkg/settings"
	"time"
)

var login = login2.NewLoginInterFace()

// GenCaptchaController 生成验证码
func GenCaptchaController(c *gin.Context) {
	time.Sleep(time.Second * 5)
	log.Logger.Infof("start")
	// 生成Token
	jwtToken, err := jwt.GenerateToken("changwenjie", "shinemo123")
	if err != nil {
		e.Error(c, e.TokenGenerateFailure)
	}
	func() {
		c.SetCookie(settings.OPS_PILOT_JWT_TOKEN, jwtToken, 60*60*24, "/", "localhost", false, false)
		c.SetCookie(settings.OPS_PILOT_USERNAME, "changwenjie", 60*60*24, "/", "localhost", false, false)
	}()
	e.Success(c, login.CaptchaGenerate())
}

func GetCookieController(c *gin.Context) {

	token, err := c.Cookie("ops_pilot_jwt_token")
	if err != nil {
		panic(e.TokenIsNotExist)
	}
	e.Success(c, token)
}

// ValidateCodeLoginController 验证码登录接口
func ValidateCodeLoginController(c *gin.Context) {
	var (
		loginRequest request.LoginRequest
	)
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		e.Error(c, e.InvalidParams)
	}

	// 生成Token
	jwtToken, err := jwt.GenerateToken(loginRequest.Username, loginRequest.Password)
	if err != nil {
		panic(e.TokenGenerateFailure)
	}
	// 设置cookie
	func() {
		c.SetCookie(settings.OPS_PILOT_JWT_TOKEN, jwtToken, 60*60*24, "/", "localhost", false, false)
		c.SetCookie(settings.OPS_PILOT_USERNAME, loginRequest.Username, 60*60*24, "/", "localhost", false, false)
	}()
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
