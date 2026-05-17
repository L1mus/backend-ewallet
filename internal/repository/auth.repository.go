package repository

import (
	"context"
	"log"
	"strings"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(ctx context.Context, full_name, email, hash_password string) (model.User, error) {
	sql := `INSERT INTO users (full_name,email,hash_password) VALUES($1,$2,$3) RETURNING id, full_name, email, created_at`
	args := []any{full_name, email, hash_password}
	var user model.User
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.FullName, &user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) CheckEmailExist(ctx context.Context, email string) error {
	sql := `SELECT email FROM users WHERE email=$1 AND deleted_at = NULL`
	args := []any{email}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		log.Println("Error :", err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			return appError.EmailAlreadyExists
		}
	}
	return nil
}
