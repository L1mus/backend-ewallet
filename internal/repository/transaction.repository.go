package repository

import (
	"context"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) GetUserPin(ctx context.Context, dbtx DBTX, userID int) (*string, error) {
	var hashPin *string
	sql := `SELECT hash_pin FROM users WHERE id = $1 AND deleted_at IS NULL`
	if err := dbtx.QueryRow(ctx, sql, userID).Scan(&hashPin); err != nil {
		return nil, err
	}
	return hashPin, nil
}

func (r *TransactionRepository) CheckReceiverExists(ctx context.Context, dbtx DBTX, receiverID int) (bool, error) {
	var count int
	sql := `SELECT COUNT(1) FROM users WHERE id = $1 AND deleted_at IS NULL`
	if err := dbtx.QueryRow(ctx, sql, receiverID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TransactionRepository) GetWalletBalance(ctx context.Context, dbtx DBTX, userID int) (float64, error) {
	var balance float64
	sql := `SELECT balance FROM wallet WHERE user_id = $1`
	if err := dbtx.QueryRow(ctx, sql, userID).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *TransactionRepository) InsertTransaction(ctx context.Context, dbtx DBTX, userID int, amount float64, txType string, activityType string) (int, error) {
	var transactionID int
	sql := `INSERT INTO transactions (user_id, amount, type, activity_type, status)
            VALUES ($1, $2, $3, $4, 'success') RETURNING id`
	if err := dbtx.QueryRow(ctx, sql, userID, amount, txType, activityType).Scan(&transactionID); err != nil {
		return 0, err
	}
	return transactionID, nil
}

func (r *TransactionRepository) InsertTransferDetail(ctx context.Context, dbtx DBTX, transactionID, receiverID int, description string) error {
	sql := `INSERT INTO transfer_details (transaction_id, receiver_id, description) VALUES ($1, $2, $3)`
	_, err := dbtx.Exec(ctx, sql, transactionID, receiverID, description)
	return err
}

func (r *TransactionRepository) ReduceWalletSender(ctx context.Context, dbtx DBTX, userID int, amount float64) error {
	sql := `UPDATE wallet SET balance = balance - $1, updated_at = NOW()
            WHERE user_id = $2 AND balance >= $1`
	result, err := dbtx.Exec(ctx, sql, amount, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return appError.InsufficientBalance
	}
	return nil
}

func (r *TransactionRepository) AddWalletReceiver(ctx context.Context, dbtx DBTX, userID int, amount float64) error {
	sql := `UPDATE wallet SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2`
	_, err := dbtx.Exec(ctx, sql, amount, userID)
	return err
}
