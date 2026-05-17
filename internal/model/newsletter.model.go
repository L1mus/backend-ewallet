package model

import "time"

type Newsletter struct {
	Id        int             `json:"id" db:"id"`
	UserID    *int            `json:"user_id" db:"user_id"`
	Email     string          `json:"email" db:"email"`
	Status    SubscribeStatus `json:"status" db:"status"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
