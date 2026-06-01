package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/cache"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	rdb            *redis.Client
}

func NewAuthService(authRepository *repository.AuthRepository, rdb *redis.Client) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		rdb:            rdb,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterDTO, error) {
	var hc pkg.HashConfig
	hc.UseRecommended()

	exists, err := s.authRepository.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return dto.RegisterDTO{}, err
	}
	if exists {
		return dto.RegisterDTO{}, appError.EmailAlreadyExists
	}

	hashPassword := hc.GenHash(req.Password)
	newUser, err := s.authRepository.Register(ctx, req.FullName, req.Email, hashPassword)
	if err != nil {
		return dto.RegisterDTO{}, err
	}
	return dto.RegisterDTO{
		Id:        newUser.Id,
		FullName:  newUser.FullName,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginDTO, error) {
	data, err := s.authRepository.Login(ctx, req.Email)
	if err != nil {
		return dto.LoginDTO{}, appError.EmailOrPassWrong
	}
	var hc pkg.HashConfig
	if err := hc.Compare(req.Password, data.HashPassword); err != nil {
		return dto.LoginDTO{}, appError.EmailOrPassWrong
	}
	hasPin := true
	if data.HashPin == nil {
		hasPin = false
	}
	claims := pkg.NewClaims(data.Id, data.FullName)
	token, _ := claims.GenJWT()
	return dto.LoginDTO{
		FullName: data.FullName,
		HasPin:   hasPin,
		Token:    token,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, claims pkg.Claims, token string) error {
	if claims.ExpiresAt == nil {
		return appError.TokenDoesntExpired
	}
	expirationTime := claims.ExpiresAt.Time

	ttl := time.Until(expirationTime)
	err := cache.SaveToBlacklist(ctx, s.rdb, token, ttl)
	if err != nil {
		return appError.InvalidateSession
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (dto.ForgotPasswordDTO, error) {
	/*
		validasi email user
		Generate token menggunakan package crypto/rand(lebih aman) dibanding math/rand
		buat waktu expirednya
		Simpan token yang sudah di buat ke DB
		simulasi token hanya dikirim via reponse bukan ke emailnya langsung
		*karna belum menerapkan package pengiriman email seperti SMTP atau library pihak ke 3 seperti gomail
	*/
	user, err := s.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.ForgotPasswordDTO{}, appError.UserNotFound
		}
		return dto.ForgotPasswordDTO{}, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return dto.ForgotPasswordDTO{}, err
	}
	token := hex.EncodeToString(tokenBytes)

	expiredAt := time.Now().Add(1 * time.Hour)

	if err := s.authRepository.InsertForgotPasswordToken(ctx, user.Id, token, expiredAt); err != nil {
		return dto.ForgotPasswordDTO{}, err
	}

	return dto.ForgotPasswordDTO{
		Token:     token,
		ExpiredAt: expiredAt,
	}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	/*
		Ambil token dari DB
		Validasi token
		Hash password baru
		Update password
		beri tanda token tersebut sudah di gunakan

	*/
	fp, err := s.authRepository.GetForgotPasswordToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appError.ForgotPasswordTokenNotFound
		}
		return err
	}

	if fp.IsUsed {
		return appError.ForgotPasswordTokenUsed
	}
	if time.Now().After(fp.ExpiredAt) {
		return appError.ForgotPasswordTokenExpired
	}

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashPassword := hc.GenHash(req.NewPassword)

	if err := s.authRepository.UpdatePassword(ctx, fp.UserID, hashPassword); err != nil {
		return err
	}
	if err := s.authRepository.MarkTokenAsUsed(ctx, fp.Id); err != nil {
		return err
	}

	return nil
}
