package v1

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/common/request"
	"opsPilot/internal/ops-monitor"
	"opsPilot/internal/pkg/e"
)

// CheckThirdPortController 检查三方端口连通性
func CheckThirdPortController(c *gin.Context) {
	var (
		req request.CheckThirdPortRequest
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	if ops_monitor.CheckThirdPort(req.Host, req.Port) {
		e.IsSuccess(c)
	}
}
