package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type ImageCategoryRepository interface {
	AddImageCategory(imageID uint, categoryID uint) error
	GetCategoriesByImageID(imageID uint) ([]models.Category, error)
}

type imageCategoryRepository struct {
	db *gorm.DB
}

func NewImageCategoryRepository(db *gorm.DB) ImageCategoryRepository {
	return &imageCategoryRepository{db: db}
}

// AddImageCategory - 이미지와 카테고리 관계 추가
func (r *imageCategoryRepository) AddImageCategory(imageID uint, categoryID uint) error {
	imageCategory := models.ImageCategory{
		ImageID:    imageID,
		CategoryID: categoryID,
	}
	return r.db.Create(&imageCategory).Error
}

// GetCategoriesByImageID - 특정 이미지에 연결된 모든 카테고리 조회
func (r *imageCategoryRepository) GetCategoriesByImageID(imageID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Table("categories").
		Select("categories.*").
		Joins("join image_categories on image_categories.category_id = categories.id").
		Where("image_categories.image_id = ?", imageID).
		Find(&categories).Error
	return categories, err
}
