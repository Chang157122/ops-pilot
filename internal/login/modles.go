package login

import (
	"fmt"
	"opsPilot/internal/pkg/db"
)

var sqlLite = db.SqlLiteDB

// 获取认证密钥
func getAuthSecretModel(username string) (secret string) {
	var (
		sql = fmt.Sprintf(`SELECT secret FROM t_auth_login WHERE username = '%s'`, username)
	)
	err := sqlLite.QueryRow(sql).Scan(&secret)
	if err != nil {
		panic(err)
	}
	return
}
