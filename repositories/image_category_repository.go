package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type imageCategoryRepository struct {
	db *gorm.DB
}

func NewImageCategoryRepository(db *gorm.DB) ImageCategoryRepository {
	return &imageCategoryRepository{db: db}
}

// AddImageCategory 이미지에 카테고리 추가
func (r *imageCategoryRepository) AddImageCategory(imageID uint, categoryID uint) error {
	imageCategory := models.ImageCategory{
		ImageID:    imageID,
		CategoryID: categoryID,
	}
	return r.db.Create(&imageCategory).Error
}

// GetCategoriesByImageID 특정 이미지에 연결된 모든 카테고리 조회
func (r *imageCategoryRepository) GetCategoriesByImageID(imageID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Table("categories").
		Select("categories.*").
		Joins("join image_categories on image_categories.category_id = categories.id").
		Where("image_categories.image_id = ?", imageID).
		Find(&categories).Error
	return categories, err
}

// AddCategoryToImage 이미지에 카테고리 추가
func (r *imageCategoryRepository) AddCategoryToImage(imageID, categoryID uint) error {
	return r.db.Exec(
		"INSERT INTO image_categories (image_id, category_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE image_id=image_id",
		imageID, categoryID,
	).Error
}

// RemoveCategoryFromImage 이미지에서 카테고리 제거
func (r *imageCategoryRepository) RemoveCategoryFromImage(imageID, categoryID uint) error {
	return r.db.Where("image_id = ? AND category_id = ?", imageID, categoryID).
		Delete(&models.ImageCategory{}).Error
}
