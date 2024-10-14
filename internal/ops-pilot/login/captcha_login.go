package login

import (
	"errors"
	"github.com/mojocn/base64Captcha"
	"opsPilot/internal/common/dao"
	"opsPilot/internal/common/request"
	"opsPilot/internal/common/resp"
	"opsPilot/internal/ops-pilot/models"
	"opsPilot/internal/pkg/authenticator"
	"opsPilot/internal/pkg/e"
	"opsPilot/internal/pkg/log"
)

var (
	digitDriver = base64Captcha.DriverDigit{
		Height:   50,
		Width:    200,
		Length:   4,   // 验证码长度
		MaxSkew:  0.7, // 倾斜
		DotCount: 2,   // 北京的点数越大,字体越模糊
	}
	// 使用内存驱动
	store = base64Captcha.DefaultMemStore
)

type CaptchaLogin struct{}

// RegisterUser 注册用户
func (l *CaptchaLogin) RegisterUser(req request.RegisterUserRequest) string {

	// 生成谷歌密钥
	secret := authenticator.GenerateSecret()

	if !models.AddLoginUser(dao.LoginDAO{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
		Secret:   secret,
	}) {
		log.Logger.Errorf("add login user failure")
		panic(e.RegisterUserFailure)
	}
	// 返回谷歌密钥,用于绑定
	return secret
}

// CaptchaGenerate 生成验证码
func (l *CaptchaLogin) CaptchaGenerate() resp.CaptchaResp {
	var ret resp.CaptchaResp
	c := base64Captcha.NewCaptcha(&digitDriver, store)
	_, b64s, _, err := c.Generate()
	if err != nil {
		log.Logger.Errorf("generate captcha failure: %v", err)
		panic(e.CaptchaGenerateFailure)
	}
	//ret.CaptchaId = id
	ret.Data = b64s
	return ret
}

// ValidateCodeLogin 验证码认证登录
func (l *CaptchaLogin) ValidateCodeLogin(req request.LoginRequest) string {
	if !models.CheckUserExists(req.Username) {
		panic(e.UserNotExist)
	}
	if !models.CheckUserPasswordCorrect(req.Username, req.Password) {
		panic(e.UserPasswordNotCorrect)
	}
	if !store.Verify("", req.Answer, true) {
		panic(e.CaptchaVerifyFailure)
	}

	return "success"
}

// GetAuthenticator 获取认证密钥
func (l *CaptchaLogin) GetAuthenticator(req request.LoginRequest) string {
	// 获取用户信息,检查是否存在密钥
	loginDAO := models.GetLoginUser(req.Username)
	if loginDAO.Secret == "" || len(loginDAO.Secret) == 0 {
		// 密钥不存在
		log.Logger.Errorf("secret is not")
		panic(errors.New("密钥不存在"))
	}
	return loginDAO.Secret

}
