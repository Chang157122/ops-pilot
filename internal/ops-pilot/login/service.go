package login

import (
	"opsPilot/internal/common/request"
	"opsPilot/internal/common/resp"
)

type LoginInterFace interface {
	CaptchaGenerate() resp.CaptchaResp
	ValidateCodeLogin(req request.LoginRequest) string
	GetAuthenticator(req request.LoginRequest)
	RegisterUser(req request.RegisterUserRequest) string
}

func NewLoginInterFace() LoginInterFace {
	return &CaptchaLogin{}
}
