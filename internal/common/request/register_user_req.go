package request

type RegisterUserRequest struct {
	Username string `gorm:"username;unique" json:"username"`
	Password string `gorm:"password" json:"password"`
	Role     string `gorm:"role" json:"role"`
	Email    string `gorm:"email" json:"email"`
}
