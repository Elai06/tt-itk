package dto

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type LoginResponse struct {
	Token string `form:"token"`
}

type RegisterRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Username string `form:"username"`
}
