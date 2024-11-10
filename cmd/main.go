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
	categoryController := initializers.InitCategoryModule(db)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/login", authController.Login)
	}

	api := r.Group("/api", middlewares.JWTAuthMiddleware()) // JWT 인증 미들웨어를 타는 API 그룹

	userAPI := api.Group("/user")              // user용 엔드포인트
	userAPI.Use(middlewares.RequireUserRole()) // user 권한 검증 미들웨어
	{
		// 이미지 API
		imageAPI := userAPI.Group("/images")
		imageAPI.POST("", imageController.UploadImage)

		imageAPI.GET("", imageController.GetImagesByUserID)
		imageAPI.GET("/thumbnail/:imageID/", imageController.GetThumbnail)              // 썸네일 조회 엔드포인트
		imageAPI.GET("/:imageID/categories", categoryController.GetCategoriesByImageID) // 카테고리별 이미지 조회
		imageAPI.GET("/:imageID/", imageController.GetImageByID)

		imageAPI.DELETE("", imageController.DeleteAllUserImages)
		imageAPI.DELETE("/:imageID", imageController.DeleteImage)

		// 카테고리 API
		categoryAPI := userAPI.Group("/categories")

		categoryAPI.GET("/:categoryID/images", imageController.GetImagesByCategoryID)
		categoryAPI.POST("/:categoryID/images/:imageID/", categoryController.AddCategoryToImage)
		categoryAPI.DELETE("/:categoryID/images/:imageID/", categoryController.RemoveCategoryFromImage)
	}

	adminAPI := api.Group("/admin")              // admin용 엔드포인트
	adminAPI.Use(middlewares.RequireAdminRole()) // admin 권한 검증 미들웨어
	{
		// 이미지 API
		imageAPI := adminAPI.Group("/images")

		imageAPI.POST(":userID/", imageController.UploadImage)

		imageAPI.GET("", imageController.GetAllImagesByAdmin)
		imageAPI.GET("users/:userID/images", imageController.GetImagesByUserID)
		imageAPI.GET("/:imageID/", imageController.GetImageByID)

		imageAPI.DELETE("/users/:userID/images", imageController.DeleteAllUserImages)
		imageAPI.DELETE("/:imageID/", imageController.DeleteImage)
		imageAPI.GET("/:imageID/categories", categoryController.GetCategoriesByImageID) // 카테고리별 이미지 조회

		// 카테고리 API
		categoryAPI := adminAPI.Group("/categories")

		categoryAPI.GET("/:categoryID/images", imageController.GetImagesByCategoryID)
		categoryAPI.POST("/:categoryID/images/:imageID/", categoryController.AddCategoryToImage)
		categoryAPI.DELETE("/:categoryID/images/:imageID/", categoryController.RemoveCategoryFromImage)
	}

	_ = r.Run()
}
