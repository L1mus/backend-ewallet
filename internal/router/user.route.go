package router

import (
	"github.com/L1mus/backend-ewallet/internal/controller"
	"github.com/L1mus/backend-ewallet/internal/middleware"
	"github.com/L1mus/backend-ewallet/internal/repository"
	"github.com/L1mus/backend-ewallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(router *gin.Engine, db *pgxpool.Pool) {
	userRouter := router.Group("/users")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	userRouter.GET("/profile", middleware.VerifyToken, userController.GetUserProfile)
	userRouter.GET("/dashboard", middleware.VerifyToken, userController.GetUserDashboard)
	userRouter.GET("/report", middleware.VerifyToken, userController.GetTransactionReport)
	userRouter.GET("/transfer", middleware.VerifyToken, userController.FindReceiver)
	userRouter.GET("/transactions", middleware.VerifyToken, userController.GetTransactionHistory)
	userRouter.PATCH("/pin", middleware.VerifyToken, userController.EditPin)
}
