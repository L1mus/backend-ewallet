package router

import (
	"github.com/L1mus/backend-ewallet/internal/controller"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/", authController.Login)
	authRouter.POST("/register", authController.Register)
}
