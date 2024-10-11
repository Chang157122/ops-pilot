package v1

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/ops-monitor"
	"opsPilot/internal/pkg/e"
)

// MetricsController 监控指标接口
func MetricsController(c *gin.Context) {
	e.Success(c, ops_monitor.GetMetrics())
}
