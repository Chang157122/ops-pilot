package e

var messages = map[int]string{
	SUCCESS:         "success",
	FAILED:          "failed",
	InvalidParams:   "无效的参数",
	TokenIsNotExist: "token不存在",
	TokenIsInvalid:  "无效的token",
	TokenTimeout:    "token过期",
}

func GetMessage(code int) string {
	return messages[code]
}
