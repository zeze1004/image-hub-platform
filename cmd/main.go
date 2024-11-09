package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/initializers"
)

func main() {
	// DB 초기화
	db := initializers.InitDB()

	// 모듈 초기화
	authController := initializers.InitUserModule(db)
	imageController := initializers.InitImageModule(db)

	authMiddleware := initializers.InitAuthMiddleware()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		// 이미지 업로드 API
		api.POST("/upload", imageController.UploadImage)
		api.POST("/upload/:userID", imageController.UploadImage) // ADMIN 계정이 USER의 이미지 업로드
	}

	_ = r.Run()
}
