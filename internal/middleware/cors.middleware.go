package middleware

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	origin := ctx.GetHeader("Origin")

	AllowOrigins := []string{os.Getenv("ALLOWED_ORIGINS")}
	AllowHeaders := []string{"Origin", "Content-Type", "Authorization"}
	AllowMethods := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions}

	if slices.Contains(AllowOrigins, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}

	ctx.Header("Access-Control-Allow-Headers", strings.Join(AllowHeaders, ", "))
	ctx.Header("Access-Control-Allow-Methods", strings.Join(AllowMethods, ", "))

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
