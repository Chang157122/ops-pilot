package models

import (
	"opsPilot/internal/common/dao"
	"opsPilot/internal/pkg/log"
)

type LoginModel struct {
	Model
	Username string `gorm:"username" json:"username"`
	Password string `gorm:"password" json:"password"`
	Secret   string `gorm:"secret" json:"secret"`
	JWTToken string `gorm:"jwt_token" json:"jwt_token"`
}

// GetOneSecret 获取认证
func GetOneSecret(username string) (LoginModel LoginModel) {
	db.Where("username=?", username).First(&LoginModel)
	return
}

// AddLoginUser 新增登录用户
func AddLoginUser(loginDAO dao.LoginDAO) bool {
	create := db.Create(
		&LoginModel{
			Username: loginDAO.Username,
			Password: loginDAO.Password,
			Secret:   loginDAO.Secret,
			JWTToken: loginDAO.JWTToken,
		},
	)
	if create.Error != nil {
		log.Logger.Errorf("Add login user failed err: %v", create.Error)
		return false
	}
	return true
}
