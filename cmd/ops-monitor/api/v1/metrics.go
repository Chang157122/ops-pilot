package v1

import (
	"github.com/gin-gonic/gin"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/watchdog"
)

func MetricsController(c *gin.Context) {
	e.Success(c, watchdog.GetMetrics())
}
