package model

import "time"

type CategoryPaymentMethod struct {
	Id           int       `json:"id" db:"id"`
	CategoryName string    `json:"category_name" db:"category_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type PaymentMethod struct {
	Id                int       `json:"id" db:"id"`
	PaymentCategoryID int       `json:"payment_category_id" db:"payment_category_id"`
	Name              string    `json:"name" db:"name"`
	Code              string    `json:"code" db:"code"`
	Fee               float64   `json:"fee" db:"fee"`
	LogoURL           *string   `json:"logo_url" db:"logo_url"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
