package dto

import (
	"time"
)

type ResponseSuccess struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Welcome, John doe"`
}

type ResponseError struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Failed get data/internal server error"`
	Error   string `json:"error" example:"internal server error/bad request"`
}

type RegisterResponse struct {
	ResponseSuccess
	Data RegisterDTO
}

type LoginResponse struct {
	ResponseSuccess
	Data LoginDTO
}

type GetUserProfileResponse struct {
	ResponseSuccess
	Data GetUserProfileDTO
}

type RegisterDTO struct {
	Id        int       `json:"id" example:"101"`
	FullName  string    `json:"full_name" example:"John Doe"`
	Email     string    `json:"email" example:"example@mail.com"`
	CreatedAt time.Time `json:"created-at" example:"2026-05-19T23:35:11.88652Z"`
}

type LoginDTO struct {
	FullName string `json:"full_name" example:"John Doe"`
	HasPin   bool   `json:"has_pin"   example:"true"`
	Token    string `json:"token"`
}

type GetUserProfileDTO struct {
	Id                int     `json:"id" example:"1"`
	FullName          string  `json:"full_name" example:"John Doe"`
	Email             string  `json:"email" example:"example@mail.com"`
	Phone             *string `json:"phone" example:"021234512552"`
	ProfilePictureURL *string `json:"profile_picture_url" example:"https://example.com"`
}

type GetUserDashboardDTO struct {
	Balance       float32 `json:"balance"`
	TotalIncome   float32 `json:"total_income"`
	TotalExpenses float32 `json:"total_expenses"`
}

type FindReceiverDTO struct {
	Id                int    `json:"id"`
	FullName          string `json:"full_name"`
	Phone             string `json:"phone"`
	ProfilePictureUrl string `json:"profile_picture_url"`
	IsVerified        bool   `json:"is_verified"`
}

type FindReceiverResponse struct {
	ResponseSuccess
	Data []FindReceiverDTO
}

type GetTransactionReportDTO struct {
	Period       string  `json:"period"`
	TotalIncome  float32 `json:"total_income"`
	TotalExpense float32 `json:"total_expense"`
}

type ReportSummaryDTO struct {
	TotalIncome  float32 `json:"total-income"`
	TotalExpense float32 `json:"total-expense"`
	NetAmount    float32 `json:"net_amount"`
}

type GetTransactionReportResponse struct {
	Period  string
	Summary ReportSummaryDTO
	Data    []GetTransactionReportDTO
}

type PaginationMetaData struct {
	TotalPages int    `json:"total-page,omitempty"`
	TotalData  int    `json:"total_data,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	NextLink   string `json:"next_page,omitempty"`
	PrevLink   string `json:"prev_page,omitempty"`
}

type GetTransactionHistoryDTO struct {
	TransactionID     int     `json:"transaction_id"`
	Amount            float32 `json:"amount"`
	Type              string  `json:"type"`
	ActivityType      string  `json:"activity_type"`
	Status            string  `json:"status"`
	ReceiverName      string  `json:"receiver_name"`
	Phone             string  `json:"phone"`
	ProfilePictureUrl string  `json:"profile_picture_url"`
}

type NewBalanceDTO struct {
	Balance float64 `json:"balance"`
}

type TopupDetailDTO struct {
	TransactionID int     `json:"transaction_id"`
	PaymentMethod string  `json:"payment_method"`
	OrderAmount   float64 `json:"order_amount"`
	Fee           float64 `json:"fee"`
	TaxAmount     float64 `json:"tax_amount"`
	TotalAmount   float64 `json:"total_amount"`
	NewBalance    float64 `json:"new_balance"`
}

type ForgotPasswordDTO struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}
