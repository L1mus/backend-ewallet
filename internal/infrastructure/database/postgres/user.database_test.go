package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

const (
	findByEmailQuery  = `SELECT id, full_name, email, hash_password`
	insertUserQuery   = `INSERT INTO users(full_name,email,hash_password) VALUES($1,$2,$3) RETURNING id`
	insertWalletQuery = `INSERT INTO wallet (user_id) VALUES ($1)`
)

const fakeUserID = "01911f6e-7b3a-7c3e-8f2a-abcdef123456"

func TestUserRepository_FindByEmail_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "full_name", "email", "hash_password", "hash_pin", "phone", "profile_picture_url", "created_at",
	}).AddRow(fakeUserID, "Limus", "limus@example.com", "hashed-pw", nil, nil, nil, now)

	mock.ExpectQuery(regexp.QuoteMeta(findByEmailQuery)).
		WithArgs("limus@example.com").
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), "limus@example.com")
	if err != nil {
		t.Fatalf("error tidak diharapkan: %v", err)
	}
	if user == nil {
		t.Fatal("expected user tidak nil")
	}
	if user.ID != fakeUserID {
		t.Errorf("expected ID %s, got %s", fakeUserID, user.ID)
	}
	if user.Email != "limus@example.com" {
		t.Errorf("expected email limus@example.com, dapat %s", user.Email)
	}
	if user.Phone != "" {
		t.Errorf("expected Phone kosong untuk NULL, dapat %q", user.Phone)
	}
	if user.ProfilePictureURL != "" {
		t.Errorf("expected ProfilePictureURL kosong untuk NULL, dapat %q", user.ProfilePictureURL)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ada expectation yang tidak terpenuhi: %v", err)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(findByEmailQuery)).
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByEmail(context.Background(), "notfound@example.com")
	if err != nil {
		t.Fatalf("sql.ErrNoRows seharusnya tidak dianggap error aplikasi, dapat: %v", err)
	}
	if user != nil {
		t.Errorf("expected user nil saat tidak ketemu, got %+v", user)
	}
}

func TestUserRepository_FindByEmail_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	dbErr := errors.New("connection reset by peer")
	mock.ExpectQuery(regexp.QuoteMeta(findByEmailQuery)).
		WithArgs("err@example.com").
		WillReturnError(dbErr)

	_, err = repo.FindByEmail(context.Background(), "err@example.com")
	if err == nil {
		t.Fatal("expected error saat koneksi/DB bermasalah, dapat nil")
	}
}

func TestUserRepository_Save_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	u := &domainUser.User{
		FullName:     "Limus",
		Email:        "limus@example.com",
		HashPassword: "hashed-pw",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WithArgs(u.FullName, u.Email, u.HashPassword).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fakeUserID))
	mock.ExpectExec(regexp.QuoteMeta(insertWalletQuery)).
		WithArgs(fakeUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := repo.Save(context.Background(), u); err != nil {
		t.Fatalf("error tidak diharapkan: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ada expectation yang tidak terpenuhi (kemungkinan Rollback tak terduga terpanggil): %v", err)
	}
}

func TestUserRepository_Save_BeginTxFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	beginErr := errors.New("connection refused")
	mock.ExpectBegin().WillReturnError(beginErr)

	u := &domainUser.User{FullName: "Limus", Email: "limus@example.com", HashPassword: "x"}
	err = repo.Save(context.Background(), u)
	if err == nil {
		t.Fatal("expected error saat BeginTx gagal, dapat nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ada expectation yang tidak terpenuhi: %v", err)
	}
}

func TestUserRepository_Save_InsertUserFails_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	insertErr := errors.New("duplicate key value violates unique constraint")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WillReturnError(insertErr)
	mock.ExpectRollback()

	u := &domainUser.User{FullName: "Limus", Email: "limus@example.com", HashPassword: "x"}
	err = repo.Save(context.Background(), u)
	if err == nil {
		t.Fatal("expected error saat insert user gagal, dapat nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("rollback seharusnya terpanggil tapi tidak terdeteksi: %v", err)
	}
}

func TestUserRepository_Save_InsertWalletFails_Rollback(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	walletErr := errors.New("insert wallet failed")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fakeUserID))
	mock.ExpectExec(regexp.QuoteMeta(insertWalletQuery)).
		WithArgs(fakeUserID).
		WillReturnError(walletErr)
	mock.ExpectRollback()

	u := &domainUser.User{FullName: "Limus", Email: "limus@example.com", HashPassword: "x"}
	err = repo.Save(context.Background(), u)
	if err == nil {
		t.Fatal("expected error saat insert wallet gagal, dapat nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("rollback seharusnya terpanggil tapi tidak terdeteksi: %v", err)
	}
}

func TestUserRepository_Save_CommitFails_NoFatalCrash(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("gagal membuat sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	commitErr := errors.New("connection lost before commit")
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fakeUserID))
	mock.ExpectExec(regexp.QuoteMeta(insertWalletQuery)).
		WithArgs(fakeUserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(commitErr)
	mock.ExpectRollback()

	u := &domainUser.User{FullName: "Limus", Email: "limus@example.com", HashPassword: "x"}
	err = repo.Save(context.Background(), u)

	if err == nil {
		t.Fatal("expected error saat Commit gagal, dapat nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ada expectation yang tidak terpenuhi: %v", err)
	}
}
