package models

import (
	"opsPilot/internal/common/dao"
	"opsPilot/internal/common/util"
	"opsPilot/internal/pkg/log"
)

type AuthLogin struct {
	Model
	Username string `gorm:"username;unique" json:"username"`
	Password string `gorm:"password" json:"password"`
	Role     string `gorm:"role" json:"role"`
	Email    string `gorm:"email" json:"email"`
	Secret   string `gorm:"secret" json:"secret"`
	JWTToken string `gorm:"jwt_token" json:"jwt_token"`
}

// GetOneSecret 获取认证
func GetOneSecret(username string) (LoginModel AuthLogin) {
	db.Where("username=?", username).First(&LoginModel)
	return
}

// CheckUserExists 检查用户是否存在
func CheckUserExists(username, passwd string) bool {
	err := db.Where("username=? AND password=?", username, util.StrToBase64(passwd))
	if err.Error != nil {
		log.Logger.Errorf("Check user exists failed err: %v", err)
		return false
	}
	return true
}

// AddLoginUser 新增登录用户
func AddLoginUser(loginDAO dao.LoginDAO) bool {
	create := db.Create(
		&AuthLogin{
			Username: loginDAO.Username,
			Password: util.StrToBase64(loginDAO.Password),
			Email:    loginDAO.Email,
			Role:     loginDAO.Role,
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
