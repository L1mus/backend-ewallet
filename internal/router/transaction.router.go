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

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	transactionRouter := router.Group("/transactions")
	transactionRouter.Use(middleware.VerifyToken, middleware.CheckBlacklist(rdb))

	transactionRepo := repository.NewTransactionRepository()
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter.POST("/transfer", transactionController.CreateTransfer)
}
