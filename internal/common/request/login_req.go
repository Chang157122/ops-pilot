package request

type LoginRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ValidateCode string `json:"validate_code"`
}

type AuthenticatorRequest struct {
	ValidateCode string `json:"validate_code"`
}
