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

	// TODO: 이미지와 카테고리 모듈 및 Auth 미들웨어 초기화
	//authMiddleware := initializers.InitAuthMiddleware()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	//api := r.Group("/api")
	//api.Use(authMiddleware)
	//{
	//	// 이미지와 카테고리 관련 엔드포인트들 추가
	//}

	_ = r.Run()
}
