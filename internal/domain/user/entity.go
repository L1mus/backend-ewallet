package user

import "time"

type User struct {
	ID           int
	FullName     string
	Email        string
	HashPassword string
	HashPin      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

type Wallet struct {
	ID            int
	UserID        int
	BalanceOnCent float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
