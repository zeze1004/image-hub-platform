package initializers

import (
	"github.com/zeze1004/image-hub-platform/controllers"
	"github.com/zeze1004/image-hub-platform/repositories"
	"github.com/zeze1004/image-hub-platform/services"
	"gorm.io/gorm"
)

func InitImageModule(db *gorm.DB) *controllers.ImageController {
	imageRepo := repositories.NewImageRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	imageCategoryRepo := repositories.NewImageCategoryRepository(db)
	imageService := services.NewImageService(imageRepo, categoryRepo, imageCategoryRepo)
	imageController := controllers.NewImageController(imageService)
	return imageController
}
