package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-cs/internal/conf"
	"go-cs/internal/server/auth/server3auth/server3authfunc"
	"net/http"
	"strings"
)

const (
	userKey = "_userKey"
)

func GetUserIdFromCtx(ctx context.Context) int64 {
	return ctx.Value(userKey).(int64)
}
func authMiddleware2Gin(confJwt *conf.Jwt, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request

		token, err := server3authfunc.AuthTokenFromHttpReq(r, confJwt, rdb)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(userKey, token.UserId)
	}
}

// cors
func Cors(conf *conf.Server) gin.HandlerFunc {
	allowOrigins := []string{"*"}
	if len(conf.Http.CorsOrigins) > 0 {
		allowOrigins = conf.Http.CorsOrigins
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", strings.Join(allowOrigins, ","))
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}
