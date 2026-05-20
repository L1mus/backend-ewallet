package repository

import (
	"context"

	"github.com/L1mus/backend-ewallet/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserProfile(ctx context.Context, id int) (model.UserProfile, error) {
	sql := `SELECT id,full_name,email,phone,profile_picture_url FROM users WHERE id = $1`

	args := []any{id}
	var data model.UserProfile
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&data.Id, &data.FullName, &data.Email, &data.Phone, &data.ProfilePictureURL); err != nil {
		return model.UserProfile{}, err
	}
	return data, nil
}

func (r *UserRepository) GetUserDashboard(ctx context.Context, id int) (model.UserDashboard, error) {
	sql := `
			SELECT  w.balance, SUM(CASE WHEN t.type = 'income' AND t.status = 'success' THEN t.amount ELSE 0 END) AS total_income, SUM(CASE WHEN t.type = 'expense' AND t.status = 'success' THEN t.amount ELSE 0 END) AS total_expense
			FROM wallet w
			LEFT JOIN transactions t ON w.user_id = t.user_id
			WHERE w.user_id = $1
			GROUP BY w.balance`
	args := []any{id}
	var data model.UserDashboard
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&data.Balance, &data.TotalIncome, &data.TotalExpenses); err != nil {
		return model.UserDashboard{}, err
	}
	return data, nil
}