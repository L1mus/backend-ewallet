package usecase

import (
	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

type AuthUseCase struct {
	userRepo domainUser.UserRepository
}

func NewAuthUseCase(repo domainUser.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: repo,
	}
}

func (c *AuthUseCase) Register() {

}
