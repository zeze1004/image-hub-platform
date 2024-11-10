package initializers

import (
	"github.com/zeze1004/image-hub-platform/controllers"
	"github.com/zeze1004/image-hub-platform/repositories"
	"github.com/zeze1004/image-hub-platform/services"
	"gorm.io/gorm"
)

func InitCategoryModule(db *gorm.DB) *controllers.CategoryController {
	imageRepo := repositories.NewImageRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	imageCategoryRepo := repositories.NewImageCategoryRepository(db)
	categoryService := services.NewCategoryService(imageRepo, categoryRepo, imageCategoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)
	return categoryController
}
