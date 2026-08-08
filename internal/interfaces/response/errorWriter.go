package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	apperror "github.com/L1mus/backend-ewallet/internal/AppError"
	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
)

func mapDomainError(err error) *apperror.AppError {
	switch {
	case errors.Is(err, domainUser.ErrEmailAlreadyExist):
		return apperror.NewFromError(http.StatusBadRequest, "EMAIL_EXIST", "Email already registered", err)
	default:
		return nil
	}
}

func WriteError(w http.ResponseWriter, err error) {
	var appError *apperror.AppError

	if errors.As(err, &appError) {
		writeJSON(w, appError)
		return
	}

	if mapped := mapDomainError(err); mapped != nil {
		writeJSON(w, mapped)
		return
	}

	log.Printf("[UNEXPECTED ERROR] %v", err)
	writeJSON(w, apperror.NewAppError(http.StatusInternalServerError, "server error"))
}

func writeJSON(w http.ResponseWriter, appError *apperror.AppError) {
	if appError.Err != nil {
		log.Printf("[ERROR] code=%s status=%d cause=%v", appError.Code, appError.StatusCode, appError.Err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.StatusCode)
	_ = json.NewEncoder(w).Encode(appError)
}
