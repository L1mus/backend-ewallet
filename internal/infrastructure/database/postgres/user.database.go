package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	SELECT id, full_name, email, hash_password, hash_pin, phone, profile_picture_url, created_at
	FROM users 
	WHERE email = $1 AND deleted_at IS NULL`

	var u domainUser.User
	var hashPin, phone, profilePictureURL sql.NullString
	row := r.db.QueryRowContext(ctx, q, email)

	err := row.Scan(
		&u.ID,
		&u.FullName,
		&u.Email,
		&u.HashPassword,
		&hashPin,
		&phone,
		&profilePictureURL,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.HashPin = hashPin.String
	u.Phone = phone.String
	u.ProfilePictureURL = profilePictureURL.String

	return &u, nil
}

func (r *userRepository) Save(ctx context.Context, u *domainUser.User) error {
	return nil
}
