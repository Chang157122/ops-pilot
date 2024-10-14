package middleeware

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/pkg/jwt"
	"opsPilot/internal/pkg/log"
	"opsPilot/internal/pkg/settings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		token, err := c.Cookie(settings.OPS_PILOT_JWT_TOKEN)
		if err != nil {
			code = e.FAILED
			log.Logger.Errorf("get Cooike failed! key: %s err: %v", settings.OPS_PILOT_JWT_TOKEN, err)
			e.IsError(c, code, err.Error())
			c.Abort()
			return
		}
		if token == "" {
			code = e.TokenIsNotExist
		}
		parseToken, err := jwt.ParseToken(token)
		if err != nil {
			code = e.TokenIsInvalid
		} else if time.Now().Unix() > parseToken.ExpiresAt {
			code = e.TokenTimeout
		}
		if code != e.SUCCESS {
			e.Error(c, code)
			c.Abort()
			return
		}
		c.Next()
	}
}
