package models

type UserRegisteration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AuthType string `json:"auth_type"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
