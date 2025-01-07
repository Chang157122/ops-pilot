package models

import (
	"opsPilot/internal/common/dao"
	"testing"
)

func TestModels(t *testing.T) {
	var loginDAO = dao.LoginDAO{
		Username: "12312a31",
		Password: "123456",
		Secret:   "123456",
	}

	AddLoginUser(loginDAO)
	//model := GetOneSecret("1231231")
	//fmt.Println(model)
}
