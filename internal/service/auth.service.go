package service

import (
	"context"
	"time"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/cache"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
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
