package dto

type RegisterRequest struct {
	FullName        string `json:"full_name" binding:"required,min=3"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=8"`
}
