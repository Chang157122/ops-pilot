package login

import "opsPilot/internal/common/request"

func Login() {

}

// GetAuthenticator 获取认证密钥
func GetAuthenticator(req request.LoginRequest) {
	secret := getAuthSecretModel(req.Username)
	if secret == "" || len(secret) == 0 {

	}
}
