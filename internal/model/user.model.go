package model

import "time"

type User struct {
	Id                int        `json:"id" db:"id"`
	FullName          string     `json:"full_name" db:"full_name"`
	Email             string     `json:"email" db:"email"`
	HashPassword      string     `json:"hash_password" db:"hash_password"`
	HashPin           *string    `json:"hash_pin" db:"hash_pin"`
	Phone             *string    `json:"phone" db:"phone"`
	ProfilePictureURL *string    `json:"profile_picture_url" db:"profile_picture_url"`
	IsVerified        bool       `json:"is_verified" db:"is_verified"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
}

type OauthUser struct {
	Id             int           `json:"id" db:"id"`
	UserID         int           `json:"user_id" db:"user_id"`
	ProviderName   OauthProvider `json:"provider_name" db:"provider_name"`
	ProviderUserID string        `json:"provider_user_id" db:"provider_user_id"`
	AccessToken    *string       `json:"access_token" db:"access_token"`
	RefreshToken   *string       `json:"refresh_token" db:"refresh_token"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	ExpiredAt      *time.Time    `json:"expired_at" db:"expired_at"`
	UpdatedAt      *time.Time    `json:"updated_at" db:"updated_at"`
}

type FavoriteContact struct {
	Id             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	FavoriteUserID int       `json:"favorite_user_id" db:"favorite_user_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type UserProfile struct {
	Id                int     `json:"id" db:"id"`
	FullName          string  `json:"full_name" db:"full_name"`
	Email             string  `json:"email" db:"email"`
	Phone             *string `json:"phone" db:"phone"`
	ProfilePictureURL *string `json:"profile_picture_url" db:"profile_picture_url"`
}

type UserDashboard struct {
	Balance       float32 `db:"balance"`
	TotalIncome   float32 `db:"total_income"`
	TotalExpenses float32 `db:"total_expenses"`
}
type FindReceiver struct {
	Id                int    `db:"id"`
	FullName          string `db:"full_name"`
	Phone             string `db:"phone"`
	ProfilePictureUrl string `db:"profile_picture_url"`
	IsVerified        bool   `db:"is_verified"`
	TotalCount        int    `db:"total_count"`
}

type GetTransactionReport struct {
	Period       *time.Time `db:"period"`
	TotalIncome  float32    `db:"total_income"`
	TotalExpense float32    `db:"total_expense"`
}

type GetTransactionHistory struct {
	TransactionID     int     `db:"transaction_id"`
	Amount            float32 `db:"amount"`
	Type              string  `db:"type"`
	ActivityType      string  `db:"activity_type"`
	Status            string  `db:"status"`
	ReceiverName      *string `db:"receiver_name"`
	Phone             *string `db:"phone_receiver"`
	ProfilePictureUrl *string `db:"profile_picture_url"`
	TotalCount        int     `db:"total_count"`
}
