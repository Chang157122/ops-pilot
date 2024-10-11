package e

var messages = map[int]string{
	SUCCESS:                "success",
	FAILED:                 "failed",
	InvalidParams:          "无效的参数",
	TokenIsNotExist:        "Token不存在",
	TokenIsInvalid:         "无效的Token",
	TokenTimeout:           "Token过期",
	TokenGenerateFailure:   "Token生成失败",
	RegisterUserFailure:    "用户注册失败",
	UserNotExist:           "用户不存在",
	CaptchaVerifyFailure:   "验证码认证失败",
	CaptchaGenerateFailure: "验证码生成失败",
	UserPasswordNotCorrect: "用户密码不正确",
}

func GetMessage(code int) string {
	return messages[code]
}
