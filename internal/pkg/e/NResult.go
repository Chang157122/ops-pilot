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
	return
}
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, NewNResult(SUCCESS, data, GetMessage(SUCCESS)))
	return
}

func Error(c *gin.Context, code any) {
	switch code.(type) {
	case int:
		Code := code.(int)
		c.JSON(http.StatusOK, NewNResult(Code, nil, GetMessage(Code)))
		return
	case error:
		c.JSON(http.StatusOK, NewNResult(FAILED, nil, code.(error).Error()))
		return
	}
}

func IsError(c *gin.Context, code int, data string) {
	c.JSON(http.StatusOK, NewNResult(code, nil, data))
	return
}
