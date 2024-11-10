package services

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/models"
)

type AuthService interface {
	SignUp(user *models.User) error
	Login(email, password string) (string, error)
}

type ImageService interface {
	UploadImage(ctx *gin.Context, fileName, description string, userID uint, categoryNames []string) (*models.Image, error)
	GetThumbnail(imageID uint) (string, error)
	GetAllImages() ([]models.Image, error)
	GetImagesByUserID(userID uint) ([]models.Image, error)
	GetImageByID(imageID uint, userID uint, isAdmin bool) (*models.Image, error)
	DeleteImageByID(imageID uint, userID uint, isAdmin bool) error
	DeleteAllImagesByUserID(userID uint) error
	GetImagesByCategoryIDAndUserID(categoryID, userID uint, isAdmin bool) ([]models.Image, error)
}

type CategoryService interface {
	GetCategoriesByImageIDAndUserID(imageID, userID uint, isAdmin bool) ([]models.Category, error)
	AddCategoryToImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error
	RemoveCategoryFromImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error
}
