package router

import (
	"database/sql"
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/application/usecase"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/database/postgres"
	"github.com/L1mus/backend-ewallet/internal/interfaces/http/handler"
)

func AuthRoute(mux *http.ServeMux, db *sql.DB) {

	userRepo := postgres.NewUserRepository(db)
	authService := usecase.NewAuthUseCase(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux.Handle("POST /auth/register", handler.Wrap(authHandler.Register))
}
