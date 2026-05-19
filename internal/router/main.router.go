package router

import (
	"net/http"

	"github.com/L1mus/backend-ewallet/internal/dto"
	"github.com/L1mus/backend-ewallet/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	// middleware global
	router.Use(middleware.CORSMiddleware)
	// router.METHOD(endpoint, callback)
	AuthRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.ResponseError{
			Status:  "Error",
			Message: "Invalid Route",
			Error:   "route not found",
		})
	})
}
