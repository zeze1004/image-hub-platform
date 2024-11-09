package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/initializers"
)

func main() {
	// DB 초기화
	db := initializers.InitDB()

	// 사용자 모듈 초기화
	authController := initializers.InitUserModule(db)

	authMiddleware := initializers.InitAuthMiddleware()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	api := r.Group("/api")
	api.Use(authMiddleware)

	_ = r.Run()
}
