package resp

// CaptchaResp 验证码响应体
type CaptchaResp struct {
	// 验证码ID
	//CaptchaId string `json:"captcha_id"`
	// 验证码图片 Base64 编码
	Data string `json:"data"`
	// 验证码数字
	//Answer string `json:"answer"`
}
