package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code,omitempty"`
	Message    string
	Errors     map[string]string `json:"errors,omitempty"`
	Err        error             `json:"-"` // contain error internal for logging, didnt serialize to JSON
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// untuk error validation client request  400
func NewAppErrorValidate(message string, fieldErrors map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Errors:     fieldErrors,
	}
}

// untuk error umum lainnya seperti 404 ,5xx , dll
func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// untuk membungkus error internal dari domain/usecase/infrastruktur
func NewFromError(statusCode int, code, safeMessage string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    safeMessage,
		Err:        err,
	}
}
