package dao

type LoginDAO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Secret   string `json:"Secret"`
}
