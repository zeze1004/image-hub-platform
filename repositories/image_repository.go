package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type ImageRepository interface {
	CreateImage(image *models.Image) error
	GetImageByID(id uint) (*models.Image, error)
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) CreateImage(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) GetImageByID(id uint) (*models.Image, error) {
	var image models.Image
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}
