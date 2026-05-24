package router

import (
	"github.com/L1mus/backend-ewallet/internal/controller"
	"github.com/L1mus/backend-ewallet/internal/middleware"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool) {
	transactionRouter := router.Group("/transactions")

	transactionRepo := repository.NewTransactionRepository()
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter.POST("/transfer", middleware.VerifyToken, transactionController.CreateTransfer)
}
