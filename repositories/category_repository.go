package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByName(names []string) ([]models.Category, error)
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
