package utils

import (
	"fmt"
	"github.com/zeze1004/image-hub-platform/repositories"
)

func ValidateImageOwnership(imageRepo repositories.ImageRepository, imageID, userID uint) error {
	image, err := imageRepo.GetImageByID(imageID)
	if err != nil {
		return fmt.Errorf("이미지를 찾을 수 없습니다: %v", err)
	}
	if image.UserID != userID {
		return fmt.Errorf("이미지에 대한 권한이 없습니다")
	}
	return nil
}
