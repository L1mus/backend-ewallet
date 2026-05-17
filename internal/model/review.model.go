package model

import "time"

type Review struct {
	Id          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Rating      int16     `json:"rating" db:"rating"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
