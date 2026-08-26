package middleware

import (
	"fmt"
	"log"
	"net/http"

	apperror "github.com/L1mus/backend-ewallet/internal/AppError"
	"github.com/L1mus/backend-ewallet/internal/interfaces/response"
)

// Recovery menangkap panic yang tidak ter-handle (nil pointer, index out of range, dll)
// supaya server tetap hidup dan client tetap dapat response JSON yang konsisten,
// bukan connection reset / server mati.

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVERED] %v", rec)
				response.WriteError(w, apperror.NewFromError(
					http.StatusInternalServerError,
					"INTERNAL_ERROR",
					"internal server error",
					fmt.Errorf("[PANIC RECOVER] %v", rec),
				))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
