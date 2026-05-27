package middleware

import (
	"context"

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
		rawToken, _ := ctx.Get("raw_token")
		tokenStr := rawToken.(string)

		isBlacklisted, err := cache.IsBlacklisted(context.Background(), rdb, tokenStr)
		if err != nil {
			response.Error(ctx, 500, "failed to validate session status")
			ctx.Abort()
			return
		}

		if isBlacklisted {
			response.Error(ctx, 401, "Session has expired, please login again")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
