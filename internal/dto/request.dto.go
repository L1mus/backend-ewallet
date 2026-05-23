package dto

type RegisterRequest struct {
	FullName        string `json:"full_name" binding:"required,min=3" example:"John Doe"`
	Email           string `json:"email" binding:"required,email" example:"example@mail.com"`
	Password        string `json:"password" binding:"required,min=8" example:"example123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"example123"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"example@mail.com"`
	Password string `json:"password" example:"example123"`
}

type GetTransactionsReportRequest struct {
	Period string `json:"period" example:"month"`
}

type TransferQuery struct {
	Page   string `json:"page" example:"1"`
	Search string `json:"search" example:"John Doe"`
}

type ReceiverQuery struct {
	Page   string `json:"page" example:"1"`
	Search string `json:"search" example:"John Doe"`
	Page   string `form:"page" default:"1" example:"1"`
	Search string `form:"search" example:"John Doe"`
}
}
