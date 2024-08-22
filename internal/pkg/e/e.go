package e

var messages = map[int]string{
	SUCCESS:         "success",
	FAILED:          "failed",
	InvalidParams:   "无效的参数",
	TokenIsNotExist: "Token不存在",
	TokenIsInvalid:  "无效的Token",
	TokenTimeout:    "Token过期",
}

func GetMessage(code int) string {
	return messages[code]
}
