package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/initializers"
	"github.com/zeze1004/image-hub-platform/middlewares"
)

func main() {
	// DB 초기화
	db := initializers.InitDB()

	// 모듈 초기화
	authController := initializers.InitUserModule(db)
	imageController := initializers.InitImageModule(db)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	api := r.Group("/api", middlewares.JWTAuthMiddleware())
	// 사용자용 엔드포인트
	userAPI := api.Group("/user")
	userAPI.Use(middlewares.RequireUserRole()) // 사용자 권한 미들웨어
	{
		userAPI.POST("/upload", imageController.UploadImage)
		userAPI.GET("/thumbnail/:imageID/", imageController.GetThumbnail) // 썸네일 조회 엔드포인트
		userAPI.GET("/images", imageController.GetImages)
		userAPI.GET("/images/:imageID/", imageController.GetImageByID)
		// 이미지 삭제
		userAPI.DELETE("/images", imageController.DeleteAllUserImages)
		userAPI.DELETE("/images/:imageID", imageController.DeleteImage)
		userAPI.GET("/images/:imageID/categories", imageController.GetCategoriesByImageID) // 카테고리별 이미지 조회

		// 카테고리 API
		categoryAPI := userAPI.Group("/categories")
		categoryAPI.GET("/:categoryID/images", imageController.GetImagesByCategoryID)
	}

	// 관리자용 엔드포인트
	adminAPI := api.Group("/admin")
	adminAPI.Use(middlewares.RequireAdminRole()) // 관리자 권한 미들웨어
	{
		adminAPI.POST("/upload/:userID/", imageController.UploadImage)
		adminAPI.GET("/images", imageController.GetAdminImages)
		adminAPI.GET("users/:userID/images", imageController.GetAdminImages)
		adminAPI.GET("/images/:imageID/", imageController.GetAdminImageByID)
		// 이미지 삭제
		adminAPI.DELETE("/users/:userID/images", imageController.DeleteAllUserImages)
		adminAPI.DELETE("/images/:imageID/", imageController.DeleteImage)
		adminAPI.GET("/images/:imageID/categories", imageController.GetCategoriesByImageID) // 카테고리별 이미지 조회

		// 카테고리 API
		categoryAPI := adminAPI.Group("/categories")
		categoryAPI.GET("/:categoryID/images", imageController.GetImagesByCategoryID)
	}

	_ = r.Run()
}
