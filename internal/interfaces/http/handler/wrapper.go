package handler

import (
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/interfaces/response"
)

// AppHandler ,handler http.handlerfunc yang mengembalikan error
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// Wrap ,jika handler return error, otomatis diteruskan ke response.WriteError.
func Wrap(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			response.WriteError(w, err)
		}
	}
}
