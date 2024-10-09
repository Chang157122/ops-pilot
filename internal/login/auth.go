package login

import (
	"opsPilot/internal/common/request"
	"opsPilot/internal/models"
	"opsPilot/internal/pkg/log"
)

func ValidateCodeLogin(req request.LoginRequest) bool {
	if models.CheckUserExists(req.Username, req.Password) {
		return true
	}
	return false
}

// GetAuthenticator 获取认证密钥
func GetAuthenticator(req request.LoginRequest) {
	loginDAO := models.GetOneSecret(req.Username)
	if loginDAO.Secret == "" || len(loginDAO.Secret) == 0 {
		// 密钥不存在
		log.Logger.Errorf("secret is not")
		return
	}

}
