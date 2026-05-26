package service

import (
	"context"
	"log"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	db              *pgxpool.Pool
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, db *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *TransactionService) CreateTransfer(ctx context.Context, senderID int, req dto.CreateTransferRequest) error {
	/*
		VALIDASI
		cegah transfer ke diri sendiri
		verifikasi PIN
		cek receiver ada
		cek balance awal

		TRANSACTION
		mulai Begin
		defer Rollback
		catat transaksi expense untuk sender
		catat detail transfer (receiver & deskripsi)
		kurangi wallet sender
		tambah wallet receiver
		selesai Commit
	*/

	if senderID == req.ReceiverID {
		return appError.SelfTransferNotAllowed
	}

	hashPin, err := s.transactionRepo.GetUserPin(ctx, s.db, senderID)
	if err != nil {
		return err
	}
	if hashPin == nil {
		return appError.EmptyPin
	}
	var hc pkg.HashConfig
	if err := hc.Compare(req.Pin, *hashPin); err != nil {
		return appError.WrongPin
	}

	exists, err := s.transactionRepo.CheckReceiverExists(ctx, s.db, req.ReceiverID)
	if err != nil {
		return err
	}
	if !exists {
		return appError.ReceiverNotFound
	}

	balance, err := s.transactionRepo.GetWalletBalance(ctx, s.db, senderID)
	if err != nil {
		return err
	}
	if balance < req.Amount {
		return appError.InsufficientBalance
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("rollback error: ", err.Error())
		}
	}(tx, ctx)

	transactionID, err := s.transactionRepo.InsertTransaction(ctx, tx, senderID, req.Amount, "expense", "transfer")
	if err != nil {
		return err
	}

	if err := s.transactionRepo.InsertTransferDetail(ctx, tx, transactionID, req.ReceiverID, req.Description); err != nil {
		return err
	}

	if err := s.transactionRepo.ReduceWalletSender(ctx, tx, senderID, req.Amount); err != nil {
		return err
	}

	if err := s.transactionRepo.AddWalletReceiver(ctx, tx, req.ReceiverID, req.Amount); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
