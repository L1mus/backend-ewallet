package repository

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/L1mus/backend-ewallet/internal/dto"
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

func (r *UserRepository) GetReceiver(ctx context.Context, id int, req dto.ReceiverQuery) ([]model.FindReceiver, error) {
	//membuat string query dengan strings.builder
	//args sebagai nilai yang akan dimasukan ke parameterization query
	//variable count yang akan terus increment sesuai dengan panjang variable args
	var sb strings.Builder
	var args []any
	argCount := 1

	sb.WriteString(`
			SELECT id, full_name, phone, profile_picture_url, is_verified
			FROM users 
			WHERE deleted_at IS NULL
			  AND id != $1
			  `)
	args = append(args, id)
	argCount++
	if req.Search != "" {
		_, err := fmt.Fprintf(&sb, `AND (full_name ILIKE %%$%d OR phone  ILIKE %%$%d)`, argCount, argCount)
		if err != nil {
			return nil, err
		}
		args = append(args, req.Search)
		argCount++
	}
	sb.WriteString(`ORDER BY full_name ASC;`)
	limit := 10
	var offset int8
	if req.Page != "" {
		page, _ := strconv.Atoi(req.Page)
		if page < 0 {
			page = 1
			offset = int8((page - 1) * limit)
		} else {
			offset = int8((page - 1) * limit)
		}
		_, err := fmt.Fprintf(&sb, `LIMIT %d OFFSET %d`, argCount, argCount+1)
		if err != nil {
			return nil, err
		}
		args = append(args, limit, offset)

	}

	sql := sb.String()
	fmt.Println(sql)
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.FindReceiver
	row, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return []model.FindReceiver{}, err
	}
	defer row.Close()
	for row.Next() {
		var user model.FindReceiver
		if err := row.Scan(&user.Id, &user.FullName, &user.Phone, &user.IsVerified); err != nil {
			return []model.FindReceiver{}, err
		}
		data = append(data, user)
	}
	return data, nil
}
