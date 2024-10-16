package authenticator

import (
	"fmt"
	"testing"
)

func TestAuthenticator(t *testing.T) {

	fmt.Println("----------------- 生成secret -------------------")
	secret := GenerateSecret()
	fmt.Println("secret:" + secret)

	fmt.Println("----------------- 信息校验----------------------")
	var code int32
	fmt.Print("请输入Google Code:")

	fmt.Scan(&code)

	fmt.Print("输入错误，请重新输入：")

	b := VerifyCode(secret, code)
	if b {
		fmt.Println("验证成功！")
	} else {
		fmt.Println("验证失败！")
	}
}
