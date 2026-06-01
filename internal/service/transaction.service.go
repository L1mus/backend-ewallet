package service

import (
	"context"
	"errors"
	"log"
	"math"

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

func (s *TransactionService) CreateTransfer(ctx context.Context, senderID int, req dto.CreateTransferRequest) (dto.NewBalanceDTO, error) {
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
		return dto.NewBalanceDTO{}, appError.SelfTransferNotAllowed
	}

	hashPin, err := s.transactionRepo.GetUserPin(ctx, s.db, senderID)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}
	if hashPin == nil {
		return dto.NewBalanceDTO{}, appError.EmptyPin
	}
	var hc pkg.HashConfig
	if err := hc.Compare(req.Pin, *hashPin); err != nil {
		return dto.NewBalanceDTO{}, appError.WrongPin
	}

	exists, err := s.transactionRepo.CheckReceiverExists(ctx, s.db, req.ReceiverID)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}
	if !exists {
		return dto.NewBalanceDTO{}, appError.ReceiverNotFound
	}

	balance, err := s.transactionRepo.GetWalletBalance(ctx, s.db, senderID)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}
	if balance < req.Amount {
		return dto.NewBalanceDTO{}, appError.InsufficientBalance
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("rollback error: ", err.Error())
		}
	}(tx, ctx)

	senderTxID, err := s.transactionRepo.InsertTransaction(ctx, tx, senderID, req.Amount, "expense", "transfer")
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}

	receiverTxID, err := s.transactionRepo.InsertTransaction(ctx, tx, req.ReceiverID, req.Amount, "income", "transfer")
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}

	if err := s.transactionRepo.InsertTransferDetail(ctx, tx, senderTxID, req.ReceiverID, req.Description); err != nil {
		return dto.NewBalanceDTO{}, err
	}

	if err := s.transactionRepo.InsertTransferDetail(ctx, tx, receiverTxID, senderID, req.Description); err != nil {
		return dto.NewBalanceDTO{}, err
	}

	if err := s.transactionRepo.ReduceWalletSender(ctx, tx, senderID, req.Amount); err != nil {
		return dto.NewBalanceDTO{}, err
	}

	if err := s.transactionRepo.AddWalletReceiver(ctx, tx, req.ReceiverID, req.Amount); err != nil {
		return dto.NewBalanceDTO{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}

	newBalance, err := s.transactionRepo.GetWalletBalance(ctx, s.db, senderID)
	if err != nil {
		return dto.NewBalanceDTO{}, err
	}

	return dto.NewBalanceDTO{
		Balance: newBalance,
	}, nil
}

func (s *TransactionService) CreateTopup(ctx context.Context, userID int, req dto.CreateTopupRequest) (dto.TopupDetailDTO, error) {

	/*
			Ambil & validasi payment method
			Hitung biaya
			fee = biaya admin dari payment method
			taxRate = 10% dari order amount
			total = orderAmount + fee + tax
		    lakukan DB Transaction
			1 mulai
			2. Insert ke table transactions return id transactions
			3. Insert ke topup_details
			4. Tambah balance user
			5 commit
			jika fail rollback dengan defer fn
			Ambil balance terbaru
	*/

	pm, err := s.transactionRepo.GetPaymentMethod(ctx, s.db, req.PaymentMethodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.TopupDetailDTO{}, appError.PaymentMethodNotFound
		}
		return dto.TopupDetailDTO{}, err
	}
	if !pm.IsActive {
		return dto.TopupDetailDTO{}, appError.PaymentMethodNotActive
	}

	orderAmount := req.Amount
	fee := pm.Fee
	taxRate := 0.10
	taxAmount := math.Round(orderAmount*taxRate*100) / 100
	totalAmount := orderAmount + fee + taxAmount

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return dto.TopupDetailDTO{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil {
			log.Println("rollback error:", err.Error())
		}
	}(tx, ctx)

	transactionID, err := s.transactionRepo.InsertTransaction(ctx, tx, userID, orderAmount, "income", "topup")
	if err != nil {
		return dto.TopupDetailDTO{}, err
	}

	if err := s.transactionRepo.InsertTopupDetail(
		ctx, tx,
		transactionID, req.PaymentMethodID,
		orderAmount, fee, taxAmount, totalAmount,
	); err != nil {
		return dto.TopupDetailDTO{}, err
	}

	if err := s.transactionRepo.AddWalletReceiver(ctx, tx, userID, orderAmount); err != nil {
		return dto.TopupDetailDTO{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TopupDetailDTO{}, err
	}

	newBalance, err := s.transactionRepo.GetWalletBalance(ctx, s.db, userID)
	if err != nil {
		return dto.TopupDetailDTO{}, err
	}

	return dto.TopupDetailDTO{
		TransactionID: transactionID,
		PaymentMethod: pm.Name,
		OrderAmount:   orderAmount,
		Fee:           fee,
		TaxAmount:     taxAmount,
		TotalAmount:   totalAmount,
		NewBalance:    newBalance,
	}, nil
}
