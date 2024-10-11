package request

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// 验证码ID
	//CaptchaId string `json:"captcha_id"`
	// 验证码数字
	Answer string `json:"answer"`
}
