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

func UserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRouter := router.Group("/users")
	userRouter.Use(middleware.VerifyToken, middleware.CheckBlacklist(rdb))
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, rdb)
	userController := controller.NewUserController(userService)

	userRouter.GET("/profile", userController.GetUserProfile)
	userRouter.GET("/dashboard", userController.GetUserDashboard)
	userRouter.GET("/report", userController.GetTransactionReport)
	userRouter.GET("/transfer", userController.FindReceiver)
	userRouter.GET("/transactions", userController.GetTransactionHistory)
	userRouter.PATCH("/profile", userController.UpdateProfile)
	userRouter.PATCH("/pin", userController.EditPin)
	userRouter.PATCH("/password", userController.EditPassword)
}
