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
