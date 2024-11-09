package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequireUserRole - 사용자 권한이 필요
func RequireUserRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != "USER" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "가입된 유저만 접근 가능합니다"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireAdminRole - 관리자 권한이 필요
func RequireAdminRole() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != "ADMIN" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "어드민 계정만 접근 가능합니다"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
