package v1

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/ops-monitor"
	"opsPilot/internal/pkg/e"
)

func MetricsController(c *gin.Context) {
	e.Success(c, ops_monitor.GetMetrics())
}
