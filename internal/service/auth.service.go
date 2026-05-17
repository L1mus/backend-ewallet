package service

import (
	"context"
	"regexp"

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
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isMatched, _ := regexp.MatchString(emailRegex, req.Email)
	if !isMatched {
		return dto.RegisterResponse{}, appError.InvalidEmailFormat
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
