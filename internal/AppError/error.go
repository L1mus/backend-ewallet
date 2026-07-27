package apperror

import (
	"net/http"
)

type AppError struct {
	StatusCode int
	Message    string
	Errors     map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

// untuk error validation client request  400
func NewAppErrorValidate(message string, fieldErrors map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
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
