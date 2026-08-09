package handler

import (
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/application/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: uc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	return nil
}
