package postgres

import (
	"context"
	"database/sql"
	"errors"

	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domainUser.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	q := `
	SELECT full_name,email,phone,profile_picture_url,created_at 
	FROM users 
	WHERE email = $1`

	var u domainUser.User
	row := r.db.QueryRowContext(ctx, q, email)

	err := row.Scan(&u.FullName, &u.Email, u.Phone, u.ProfilePictureURL, u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Save(ctx context.Context, u *domainUser.User) error {
	return nil
}
