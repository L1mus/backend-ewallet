package middleware

import (
	"github.com/L1mus/backend-ewallet/internal/cache"
	"github.com/L1mus/backend-ewallet/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CheckBlacklist(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		/*
			ambil token mentah
			Cek ke Redis apakah token ini sudah di blacklist
		*/
		rawToken, exists := ctx.Get("raw_token")
		if !exists {
			response.Error(ctx, 401, "unauthorized access, please login")
			ctx.Abort()
			return
		}
		tokenStr, ok := rawToken.(string)
		if !ok || tokenStr == "" {
			response.Error(ctx, 401, "unauthorized access, please login")
			ctx.Abort()
			return
		}

		isBlacklisted, err := cache.IsBlacklisted(ctx.Request.Context(), rdb, tokenStr)
		if err != nil {
			response.Error(ctx, 500, "failed to validate session status")
			ctx.Abort()
			return
		}

		if isBlacklisted {
			response.Error(ctx, 401, "session has expired, please login again")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
