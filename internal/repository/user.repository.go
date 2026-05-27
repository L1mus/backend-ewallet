package repository

import (
	"context"
	"fmt"
	"log"
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
	sql := `SELECT id,full_name,email,phone,profile_picture_url FROM users WHERE id = $1 AND deleted_at IS NULL`

	args := []any{id}
	var data model.UserProfile
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&data.Id, &data.FullName, &data.Email, &data.Phone, &data.ProfilePictureURL); err != nil {
		return model.UserProfile{}, err
	}
	return data, nil
}

func (r *UserRepository) GetUserDashboard(ctx context.Context, id int) (model.UserDashboard, error) {
	sql := `
			SELECT  w.balance, COALESCE(SUM(CASE WHEN t.type = 'income' AND t.status = 'success' THEN t.amount ELSE 0 END), 0) AS total_income, COALESCE(SUM(CASE WHEN t.type = 'expense' AND t.status = 'success' THEN t.amount ELSE 0 END),0) AS total_expense
			FROM wallet w
			LEFT JOIN transactions t ON w.user_id = t.user_id
			WHERE w.user_id = $1
			GROUP BY w.user_id, w.balance`
	args := []any{id}
	var data model.UserDashboard
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&data.Balance, &data.TotalIncome, &data.TotalExpenses); err != nil {
		return model.UserDashboard{}, err
	}
	return data, nil
}

func (r *UserRepository) GetReceiver(ctx context.Context, id int, req dto.PageQuery) ([]model.FindReceiver, error) {
	//membuat string query dengan strings.builder
	//args sebagai nilai yang akan dimasukan ke parameterization query
	//variable count yang akan terus increment sesuai dengan panjang variable args
	var sb strings.Builder
	var args []any
	argCount := 1

	sb.WriteString(`
			SELECT id, full_name, COALESCE(phone, '') AS phone, COALESCE(profile_picture_url, '') AS profile_picture_url, is_verified, COUNT(*) OVER() AS total_count
			FROM users
			WHERE deleted_at IS NULL
			  AND id != $1
			  `)
	args = append(args, id)
	argCount++
	if req.Search != "" {
		_, err := fmt.Fprintf(&sb, `AND (full_name ILIKE  $%d  OR phone  ILIKE $%d )`, argCount, argCount)
		if err != nil {
			return nil, err
		}
		args = append(args, "%"+req.Search+"%")
		argCount++
	}
	sb.WriteString(`ORDER BY full_name ASC`)

	limit := 10
	page := 1
	if req.Page != "" {
		if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
			page = p
		}
	}
	offset := (page - 1) * limit

	_, err := fmt.Fprintf(&sb, ` LIMIT $%d OFFSET $%d`, argCount, argCount+1)
	if err != nil {
		return nil, err
	}
	args = append(args, limit, offset)

	sql := sb.String()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.FindReceiver
	for rows.Next() {
		var user model.FindReceiver
		if err := rows.Scan(&user.Id, &user.FullName, &user.Phone, &user.ProfilePictureUrl, &user.IsVerified, &user.TotalCount); err != nil {
			return nil, err
		}
		data = append(data, user)
	}
	return data, nil
}

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

func (r *UserRepository) UpdatePin(ctx context.Context, id int, hashPin string) error {
	sql := `UPDATE users SET hash_pin = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, sql, hashPin, id)
	return err
}

func (r *UserRepository) GetTransactionHistory(ctx context.Context, id int, req dto.PageQuery) ([]model.GetTransactionHistory, error) {
	var sb strings.Builder
	var args []any
	argCount := 1
	sb.WriteString(`
			SELECT  t.id AS transaction_id, t.amount, t.type, t.activity_type, u_receiver.full_name AS receiver_name,u_receiver.phone AS phone_receiver,u_receiver.profile_picture_url,COUNT(*) OVER() AS total_count
			FROM transactions t
			LEFT JOIN transfer_details td ON t.id = td.transaction_id
			LEFT JOIN users u_receiver ON td.receiver_id = u_receiver.id
			WHERE t.user_id = $1`)
	argCount++
	args = append(args, id)
	if req.Search != "" {
		sb.WriteString(`
					AND (
				  u_receiver.full_name ILIKE $2 OR
				u_receiver.phone ILIKE $2
			  )`)
		argCount++
		args = append(args, "%"+req.Search+"%")
	}
	sb.WriteString(` ORDER BY t.created_at DESC `)
	limit := 10
	page := 1
	if req.Page != "" {
		if p, _ := strconv.Atoi(req.Page); p > 0 {
			page = p
		}
	}
	offset := (page - 1) * limit
	_, err := fmt.Fprintf(&sb, `LIMIT $%d OFFSET $%d;`, argCount, argCount+1)
	if err != nil {
		return nil, err
	}
	args = append(args, limit, offset)

	//sql := `
	//		SELECT  t.id AS transaction_id, t.amount, t.type, t.activity_type, t.status, t.created_at, td.description AS transfer_description, u_receiver.full_name AS receiver_name, pm.name AS payment_method_name,COUNT(*) OVER() AS total_count
	//		FROM transactions t
	//		LEFT JOIN transfer_details td ON t.id = td.transaction_id
	//		LEFT JOIN users u_receiver ON td.receiver_id = u_receiver.id
	//		LEFT JOIN topup_details tp ON t.id = tp.transaction_id
	//		LEFT JOIN payment_method pm ON tp.payment_method_id = pm.id
	//		WHERE t.user_id = $1
	//		  AND (
	//			  u_receiver.full_name ILIKE '%' || $2 || '%' OR
	//			  pm.name ILIKE '%' || $2 || '%' OR
	//			  td.description ILIKE '%' || $2 || '%'
	//		  )
	//		ORDER BY t.created_at DESC
	//		LIMIT $3 OFFSET $4;`
	sql := sb.String()
	println(sql)
	var transactions []model.GetTransactionHistory
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.GetTransactionHistory
		if err := rows.Scan(&transaction.TransactionID, &transaction.Amount, &transaction.Type, &transaction.ActivityType, &transaction.ReceiverName, &transaction.Phone, &transaction.ProfilePictureUrl, &transaction.TotalCount); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	log.Println(transactions)
	return transactions, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id int, req dto.EditProfileRequest) error {
	sql := `
		UPDATE users
        SET full_name           = $1,
            phone               = COALESCE($2, phone),
            profile_picture_url = COALESCE($3, profile_picture_url),
            updated_at          = NOW()
        WHERE id = $4 AND deleted_at IS NULL`
	args := []any{req.FullName, req.Phone, req.ProfilePictureURL, id}
	_, err := r.db.Exec(ctx, sql, args...)
	return err
}

func (r *UserRepository) GetHashPassword(ctx context.Context, id int) (model.User, error) {
	sql := `SELECT hash_password FROM users WHERE id = $1 AND deleted_at IS NULL`
	var user model.User
	if err := r.db.QueryRow(ctx, sql, id).Scan(&user.HashPassword); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, hashPassword string) error {
	sql := `UPDATE users SET hash_password = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, sql, hashPassword, id)
	return err
}

func (r *UserRepository) CheckPhoneTaken(ctx context.Context, id int, phone string) (bool, error) {
	var count int
	sql := `SELECT COUNT(1) FROM users WHERE phone = $1 AND id != $2 AND deleted_at IS NULL`
	if err := r.db.QueryRow(ctx, sql, phone, id).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
