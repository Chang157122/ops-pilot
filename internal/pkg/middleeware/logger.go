package middleeware

import (
	"github.com/gin-gonic/gin"
	"net"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/pkg/log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 用于替换gin框架的Logger中间件，不传参数，直接这样写
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		log.Logger.Infof("path: %s | method: %s | status: %d | ip: %v | cost: %v | query: %s",
			path,
			c.Request.Method,
			c.Writer.Status(),
			c.ClientIP(),
			cost,
			query,
		)
	}
}

// GinRecovery 用于替换gin框架的Recovery中间件，因为传入参数，再包一层
func GinRecovery(stack bool) gin.HandlerFunc {
	//logger := zap.L()
	return func(c *gin.Context) {
		defer func() {
			// defer 延迟调用，出了异常，处理并恢复异常，记录日志
			if err := recover(); err != nil {
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------开始--------
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if brokenPipe {
					// 如果连接已断开，我们无法向其写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------结束--------

				// 是否打印堆栈信息，使用的是debug.Stack()，传入false，在日志中就没有堆栈信息
				if stack {
					log.Logger.Errorf("[Recovery from panic] error: %v stack: %v", err, string(debug.Stack()))
				} else {
					log.Logger.Errorf("[Recovery from panic] error: %v", err)

				}
				// 有错误，直接返回给前端错误，前端直接报错
				//c.AbortWithStatus(http.StatusInternalServerError)
				// 该方式前端不报错
				e.Error(c, err)
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
