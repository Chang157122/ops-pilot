package dao

type LoginDAO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"Secret"`
	JWTToken string `json:"jwt_token"`
}
