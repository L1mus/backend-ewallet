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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if committed {
			return
		}
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Printf("save user: rollback failed: %v", rbErr)
		}
	}()

	qUser := `INSERT INTO users(full_name,email,hash_password) VALUES($1,$2,$3) RETURNING id`

	var userID string

	err = tx.QueryRowContext(ctx, qUser, u.FullName, u.Email, u.HashPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	qWallet := `INSERT INTO wallet (user_id) VALUES ($1)`
	if _, err = tx.ExecContext(ctx, qWallet, userID); err != nil {
		return fmt.Errorf("insert wallet: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	committed = true

	return nil
}
