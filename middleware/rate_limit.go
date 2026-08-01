package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := "rate_limit:" + ctx.ClientIP()
		c := context.Background()

		// Incrementa o contador de requisições
		count, err := rdb.Incr(c, key).Result()
		if err != nil {
			ctx.Next()
			return
		}
		if count == 1 {
			_ = rdb.Expire(c, key, window).Err()
		}

		if count > int64(maxRequests) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "TOO_MANY_REQUESTS",
				"message": "Voce excedeu o limite de requisições, tente novamente mais tarde.",
				"limit":   strconv.Itoa(maxRequests),
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
