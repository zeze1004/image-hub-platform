package initializers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/middlewares"
)

// InitAuthMiddleware Auth 미들웨어 초기화
func InitAuthMiddleware() gin.HandlerFunc {
	return middlewares.JWTAuthMiddleware()
}
