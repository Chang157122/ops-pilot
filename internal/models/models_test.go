package models

import (
	"opsPilot/internal/common/dao"
	"testing"
)

func TestModels(t *testing.T) {
	var loginDAO = dao.LoginDAO{
		Username: "changwenjie",
		Password: "123456",
		Secret:   "123456",
		JWTToken: "123456",
	}
	CreateTable()
	AddLoginUser(loginDAO)
}
