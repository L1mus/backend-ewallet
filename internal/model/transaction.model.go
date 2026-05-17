package model

import "time"

type Transaction struct {
	Id           int                     `json:"id" db:"id"`
	UserID       int                     `json:"user_id" db:"user_id"`
	Amount       float64                 `json:"amount" db:"amount"`
	Type         TypeTransaction         `json:"type" db:"type"`
	ActivityType TypeActivityTransaction `json:"activity_type" db:"activity_type"`
	Status       StatusTransaction       `json:"status" db:"status"`
	CreatedAt    time.Time               `json:"created_at" db:"created_at"`
}

type TransferDetail struct {
	Id            int       `json:"id" db:"id"`
	TransactionID int       `json:"transaction_id" db:"transaction_id"`
	ReceiverID    int       `json:"receiver_id" db:"receiver_id"`
	Description   *string   `json:"description" db:"description"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type TopupDetail struct {
	Id              int       `json:"id" db:"id"`
	TransactionID   int       `json:"transaction_id" db:"transaction_id"`
	PaymentMethodID int       `json:"payment_method_id" db:"payment_method_id"`
	OrderAmount     float64   `json:"order_amount" db:"order_amount"`
	DeliveryFee     float64   `json:"delivery_fee" db:"delivery_fee"`
	TaxAmount       float64   `json:"tax_amount" db:"tax_amount"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
