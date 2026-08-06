package user

import "time"

type User struct {
	ID                string
	FullName          string
	Email             string
	HashPassword      string
	HashPin           string
	Phone             string
	ProfilePictureURL string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

type Wallet struct {
	UserID        string
	BalanceOnCent float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
