package dto

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserRegisterResponse struct {
	ID int64 `json:"id"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
