package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/utils"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token이 없습니다"})
			ctx.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "잘못된 token 입니다"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
