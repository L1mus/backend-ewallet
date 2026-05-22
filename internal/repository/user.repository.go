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
	for rows.Next() {
		var user model.FindReceiver
		if err := rows.Scan(&user.Id, &user.FullName, &user.Phone, &user.ProfilePictureUrl, &user.IsVerified); err != nil {
			return nil, err
		}
		data = append(data, user)
	}
	return data, nil
}

func (r *UserRepository) GetTotalPageReceiver(ctx context.Context, id int, req dto.ReceiverQuery) (int, int, error) {
	var sb strings.Builder
	var args []any
	argCount := 1

	sb.WriteString(`
			SELECT COUNT(DISTINCT id)
			FROM users
			WHERE deleted_at IS NULL
			  AND id != $1
			  `)
	args = append(args, id)
	argCount++

	if req.Search != "" {
		_, err := fmt.Fprintf(&sb, `AND (u.full_name ILIKE %%$%d OR u.phone  ILIKE %%$%d)`, argCount, argCount)
		if err != nil {
			return 0, 0, err
		}
		args = append(args, req.Search)
	}

	var totalReceiver int
	sql := sb.String()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&totalReceiver)
	if err != nil {
		return 0, 0, err
	}
	receiverPerPage := 10
	totalPage := int(math.Ceil(float64(totalReceiver) / float64(receiverPerPage)))
	return totalReceiver, totalPage, nil
}

//func (r *UserRepository) GetAllPagination(ctx context.Context, pagination int) {
//
//}

func (r *UserRepository) GetTransactionReport(ctx context.Context, id int, timePeriod string) ([]model.GetTransactionReport, error) {
	sql := `
			SELECT  DATE_TRUNC($2, created_at)::DATE AS period, COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
			FROM transactions
			WHERE user_id = $1 AND status = 'success'
			GROUP BY DATE_TRUNC($2, created_at)
			ORDER BY period ASC;`
	args := []any{id, timePeriod}
	var data []model.GetTransactionReport
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.GetTransactionReport
		if err := rows.Scan(&transaction.Period, &transaction.TotalIncome, &transaction.TotalExpense); err != nil {
			return nil, err
		}
		data = append(data, transaction)
	}
	return data, nil
}

func (r *UserRepository) GetPin(ctx context.Context, id int) (model.User, error) {
	sql := `SELECT hash_pin FROM users WHERE id = $1`
	args := []any{id}
	var pin model.User
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&pin.HashPin); err != nil {
		return model.User{}, err
	}
	return pin, nil
}

func (r *UserRepository) GetTransactionHistory(ctx context.Context, id int, search string, limit int8, offset int8) ([]model.GetTransactionHistory, error) {
	sql := `
			SELECT  t.id AS transaction_id, t.amount, t.type, t.activity_type, t.status, t.created_at, td.description AS transfer_description, u_receiver.full_name AS receiver_name, pm.name AS payment_method_name
			FROM transactions t
			LEFT JOIN transfer_details td ON t.id = td.transaction_id
			LEFT JOIN users u_receiver ON td.receiver_id = u_receiver.id
			LEFT JOIN topup_details tp ON t.id = tp.transaction_id
			LEFT JOIN payment_method pm ON tp.payment_method_id = pm.id
			WHERE t.user_id = $1
			  AND (
				  u_receiver.full_name ILIKE '%' || $2 || '%' OR
				  pm.name ILIKE '%' || $2 || '%' OR
				  td.description ILIKE '%' || $2 || '%'
			  )
			ORDER BY t.created_at DESC
			LIMIT $3 OFFSET $4;`
	args := []any{id, search, limit, offset}
	var transactions []model.GetTransactionHistory
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.GetTransactionHistory
		if err := rows.Scan(&transaction.TransactionID, &transaction.Amount, &transaction.Amount, &transaction.Type, &transaction.ActivityType, &transaction.Status, &transaction.CreatedAt, &transaction.PaymentMethodName, &transaction.ReceiverName); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
