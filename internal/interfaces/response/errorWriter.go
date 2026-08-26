package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	apperror "github.com/L1mus/backend-ewallet/internal/AppError"
	domainAuth "github.com/L1mus/backend-ewallet/internal/domain/auth"
	domainUser "github.com/L1mus/backend-ewallet/internal/domain/user"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/mail"
)

func mapDomainError(err error) *apperror.AppError {
	switch {
	// --- domain: user ---
	case errors.Is(err, domainUser.ErrEmailAlreadyExist):
		return apperror.NewFromError(http.StatusConflict, "EMAIL_EXISTS", "Email has been registered", err)
	case errors.Is(err, domainUser.ErrInvalidCredential):
		return apperror.NewFromError(http.StatusUnauthorized, "INVALID_CREDENTIAL", "Incorrect email or password", err)

	// --- domain: transaction ---

	// --- domain: auth ---
	case errors.Is(err, domainAuth.ErrInvalidVerificationCode):
		return apperror.NewFromError(http.StatusBadRequest, "INVALID_VERIFICATION_CODE", "Invalid verification code", err)

	// --- infrastructure: mail ---
	case errors.Is(err, mail.ErrMissingSMTPCredentials),
		errors.Is(err, mail.ErrBuildEmailBody),
		errors.Is(err, mail.ErrSendMail):
		return apperror.NewFromError(http.StatusInternalServerError, "MAIL_ERROR", "Failed to send verification email", err)

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
