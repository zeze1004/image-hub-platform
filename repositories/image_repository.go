package repositories

import (
	"github.com/zeze1004/image-hub-platform/models"
	"gorm.io/gorm"
)

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) CreateImageMetaData(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *imageRepository) GetImageByID(id uint) (*models.Image, error) {
	var image models.Image
	if err := r.db.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) GetImagesByUserID(userID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.db.Where("user_id = ?", userID).Find(&images).Error
	return images, err
}

// GetAllImages - 모든 이미지 목록 조회
func (r *imageRepository) GetAllImages() ([]models.Image, error) {
	var images []models.Image
	err := r.db.Find(&images).Error
	return images, err
}

// DeleteImage - 특정 이미지 삭제
func (r *imageRepository) DeleteImage(imageID uint) error {
	return r.db.Delete(&models.Image{}, imageID).Error
}

// DeleteImagesByUserID - 특정 사용자에 속하는 모든 이미지 일괄 삭제
func (r *imageRepository) DeleteImagesByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Image{}).Error
}
