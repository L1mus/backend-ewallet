package router

import "net/http"

func MainRoute(handler *http.Handler) http.Handler {
	mux := http.NewServeMux()

	//users
	mux.Handle("POST /api/auth")

	return mux
}
