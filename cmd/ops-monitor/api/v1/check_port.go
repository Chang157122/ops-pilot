package v1

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/common/request"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/watchdog"
)

func CheckThirdPortController(c *gin.Context) {
	var (
		req request.CheckThirdPortRequest
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	if watchdog.CheckThirdPort(req.Host, req.Port) {
		e.IsSuccess(c)
	}
}
