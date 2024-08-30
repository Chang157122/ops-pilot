package login

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"opsPilot/internal/pkg/e"
)

func GenerateCaptchaController(c *gin.Context) {
	d := struct {
		CaptchaId string
	}{captcha.New()}
	e.Success(c, d.CaptchaId)
}