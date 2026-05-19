package service

import (
	"context"

	"github.com/L1mus/backend-ewallet/internal/appError"
	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/pkg"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	var hc pkg.HashConfig
	hc.UseRecommended()

	exists, err := s.authRepository.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	if exists {
		return dto.RegisterResponse{}, appError.EmailAlreadyExists
	}

	hashPassword := hc.GenHash(req.Password)
	newUser, err := s.authRepository.Register(ctx, req.FullName, req.Email, hashPassword)
	if err != nil {
		return dto.RegisterResponse{}, err
	}
	return dto.RegisterResponse{
		Id:        newUser.Id,
		FullName:  newUser.FullName,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	data, err := s.authRepository.Login(ctx, req.Email)
	fmt.Println("Data Login:", data)
	if err != nil {
		return dto.LoginResponse{}, appError.EmailOrPassWrong
	}
	var hc pkg.HashConfig
	if err := hc.Compare(req.Password, data.HashPassword); err != nil {
		return dto.LoginResponse{}, appError.EmailOrPassWrong
	}
	claims := pkg.NewClaims(data.Id, data.FullName)
	token, err := claims.GenJWT()
	return dto.LoginResponse{
		FullName: data.FullName,
		Email:    data.Email,
		Token:    token,
	}, nil
}
