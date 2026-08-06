package user

import "context"

type UserRepository interface {
	FindEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, u *User) error
}
