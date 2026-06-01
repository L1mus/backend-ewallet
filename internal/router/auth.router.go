package router

import (
	"github.com/L1mus/backend-ewallet/internal/controller"
	"github.com/L1mus/backend-ewallet/internal/middleware"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, rdb)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/", authController.Login)
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/logout", middleware.VerifyToken, middleware.CheckBlacklist(rdb), authController.Logout)
	authRouter.POST("/forgot-password", authController.ForgotPassword)
	authRouter.POST("/reset-password", authController.ResetPassword)
}
