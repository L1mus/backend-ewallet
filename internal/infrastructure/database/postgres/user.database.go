package postgres

import (
	"context"
	"database/sql"

	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domainUser.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindEmail(ctx context.Context, email string) (*domainUser.User, error) {
	var user domainUser.User
	return &user, nil
}

func (r *userRepository) Save(ctx context.Context, u *domainUser.User) error {
	return nil
}
