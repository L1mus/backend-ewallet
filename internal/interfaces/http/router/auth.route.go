package router

import (
	"database/sql"
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/application/usecase"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/crypto"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/database/postgres"
	"github.com/L1mus/backend-ewallet/internal/interfaces/http/handler"
)

func AuthRoute(mux *http.ServeMux, db *sql.DB) {

	hasher := crypto.NewArgon2idHasher(64*1024, 1, 32, 16, 1)
	userRepo := postgres.NewUserRepository(db)
	authService := usecase.NewAuthUseCase(userRepo, hasher)
	authHandler := handler.NewAuthHandler(authService)

	mux.Handle("POST /auth/register", handler.Wrap(authHandler.Register))
}
