package middleware

import (
	"log"
	"net/http"
)

// Recovery menangkap panic yang tidak ter-handle (nil pointer, index out of range, dll)
// supaya server tetap hidup dan client tetap dapat response JSON yang konsisten,
// bukan connection reset / server mati.

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[PANIC RECOVER] %v", rec)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
