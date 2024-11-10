package repositories

import "github.com/zeze1004/image-hub-platform/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type ImageRepository interface {
	CreateImageMetaData(image *models.Image) error
	GetImageByID(id uint) (*models.Image, error)
	GetImagesByUserID(userID uint) ([]models.Image, error)
	GetAllImages() ([]models.Image, error)
	DeleteImage(imageID uint) error
	DeleteImagesByUserID(userID uint) error
}

type CategoryRepository interface {
	GetCategoriesByName(names []string) ([]models.Category, error)
	GetCategoriesByImageID(imageID uint) ([]models.Category, error)
	GetImagesByCategoryID(categoryID uint) ([]models.Image, error)
	GetImagesByCategoryIDAndUserID(categoryID, userID uint) ([]models.Image, error)
}

type ImageCategoryRepository interface {
	AddImageCategory(imageID uint, categoryID uint) error
	GetCategoriesByImageID(imageID uint) ([]models.Category, error)
	AddCategoryToImage(imageID, categoryID uint) error
	RemoveCategoryFromImage(imageID, categoryID uint) error
}
