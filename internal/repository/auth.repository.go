package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/L1mus/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5"
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

func (r *AuthRepository) Register(ctx context.Context, fullName, email, hashPassword string) (model.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("error: ", err.Error())
		}
	}(tx, ctx)

	sqlUser := `INSERT INTO users (full_name,email,hash_password) VALUES($1,$2,$3) RETURNING id, full_name, email, created_at`
	args := []any{fullName, email, hashPassword}
	var user model.User
	if err := tx.QueryRow(ctx, sqlUser, args...).Scan(&user.Id, &user.FullName, &user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}

	fmt.Println(user.Id)

	sqlWallet := `INSERT INTO wallet (user_id, balance, updated_at) VALUES($1, 0.00, NOW())`
	if _, err := tx.Exec(ctx, sqlWallet, user.Id); err != nil {
		return model.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	var count int
	sql := `SELECT COUNT(1) FROM users WHERE email=$1 AND deleted_at IS NULL`
	args := []any{email}
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AuthRepository) Login(ctx context.Context, email string) (model.User, error) {
	sql := `SELECT id,full_name,email, hash_password FROM users WHERE email = $1`
	args := []any{email}
	var user model.User
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.FullName, &user.Email, &user.HashPassword); err != nil {
		return model.User{}, err
	}
	return user, nil
}
