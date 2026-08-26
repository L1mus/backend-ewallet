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
	var payload dto.RegisterRequest
	if err := request.DecodeAndValidate(r, &payload); err != nil {
		return err
	}

	if err := h.authUseCase.Register(r.Context(), payload); err != nil {
		return err
	}

	response.WriteSuccess(w, http.StatusCreated, nil)
	return nil
}
