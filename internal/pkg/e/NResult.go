package e

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type NResult struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func NewNResult(code int, data any, message string) NResult {
	return NResult{Code: code, Data: data, Message: message}
}

func IsSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, NewNResult(SUCCESS, nil, GetMessage(SUCCESS)))
}
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewNResult(SUCCESS, data, GetMessage(SUCCESS)))
}

func Error(c *gin.Context, code any) {
	switch code.(type) {
	case int:
		Code := code.(int)
		c.JSON(http.StatusOK, NewNResult(Code, nil, GetMessage(Code)))
	case error:
		c.JSON(http.StatusOK, NewNResult(FAILED, nil, code.(string)))
	default:
		c.JSON(http.StatusOK, NewNResult(FAILED, nil, code.(string)))
	}
}
