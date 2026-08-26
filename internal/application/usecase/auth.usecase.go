package usecase

import (
	"context"

	"github.com/L1mus/backend-ewallet/internal/application/dto"
	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/crypto"
)

type AuthUseCase struct {
	userRepo       domainUser.UserRepository
	PasswordHasher *crypto.Argon2idHasher
}

func NewAuthUseCase(repo domainUser.UserRepository, hasher *crypto.Argon2idHasher) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       repo,
		PasswordHasher: hasher,
	}
}

func (c *AuthUseCase) Register(ctx context.Context, payload dto.RegisterRequest) error {
	isEmailExist, err := c.userRepo.FindByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}
	if isEmailExist != nil {
		return domainUser.ErrEmailAlreadyExist
	}
	hashPassword := c.PasswordHasher.Generate(payload.Password)

	newUser := &domainUser.User{
		FullName:     payload.FullName,
		Email:        payload.Email,
		HashPassword: hashPassword,
	}

	if err := c.userRepo.Save(ctx, newUser); err != nil {
		return err
	}

	return nil
}
