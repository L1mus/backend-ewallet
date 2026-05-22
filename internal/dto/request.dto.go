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

}
