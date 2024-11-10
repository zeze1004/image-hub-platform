package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByName(names []string) ([]models.Category, error)
	GetCategoriesByImageID(imageID uint) ([]models.Category, error)
	GetImagesByCategoryID(categoryID uint) ([]models.Image, error)
	GetImagesByCategoryIDAndUserID(categoryID, userID uint) ([]models.Image, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetCategoriesByName(names []string) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Where("name IN ?", names).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoriesByImageID 특정 이미지에 속한 카테고리 조회
func (r *categoryRepository) GetCategoriesByImageID(imageID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.
		Table("categories").
		Select("categories.*").
		Joins("JOIN image_categories ON image_categories.category_id = categories.id").
		Where("image_categories.image_id = ?", imageID).
		Find(&categories).Error
	return categories, err
}

// GetImagesByCategoryID 특정 카테고리에 속한 이미지 조회
func (r *categoryRepository) GetImagesByCategoryID(categoryID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.db.
		Table("images").
		Select("images.*").
		Joins("JOIN image_categories ON image_categories.image_id = images.id").
		Where("image_categories.category_id = ?", categoryID).
		Find(&images).Error
	return images, err
}

// GetImagesByCategoryIDAndUserID 특정 카테고리에 속한 사용자의 이미지 조회
func (r *categoryRepository) GetImagesByCategoryIDAndUserID(categoryID, userID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.db.
		Table("images").
		Select("images.*").
		Joins("JOIN image_categories ON image_categories.image_id = images.id").
		Where("image_categories.category_id = ? AND images.user_id = ?", categoryID, userID).
		Find(&images).Error
	return images, err
}
