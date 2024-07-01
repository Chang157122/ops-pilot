package middleeware

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/pkg/jwt"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			panic(e.TokenIsNotExist)
		}
		parseToken, err := jwt.ParseToken(token)
		if err != nil {
			panic(e.TokenIsInvalid)
		} else if time.Now().Unix() > parseToken.ExpiresAt {
			panic(e.TokenTimeout)
		}
		c.Next()
	}
}
