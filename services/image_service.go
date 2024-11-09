package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/repositories"
	"os"
	"path/filepath"
	"time"
)

type ImageService interface {
	UploadImage(ctx *gin.Context, fileName, description string, userID uint, categoryNames []string) (*models.Image, error)
}

type imageService struct {
	imageRepo         repositories.ImageRepository
	categoryRepo      repositories.CategoryRepository
	imageCategoryRepo repositories.ImageCategoryRepository
}

func NewImageService(imageRepo repositories.ImageRepository, categoryRepo repositories.CategoryRepository, imageCategoryRepo repositories.ImageCategoryRepository) ImageService {
	return &imageService{imageRepo: imageRepo, categoryRepo: categoryRepo, imageCategoryRepo: imageCategoryRepo}
}

// UploadImage - 이미지 파일을 저장하고 메타데이터를 DB에 저장
func (s *imageService) UploadImage(ctx *gin.Context, fileName, description string, userID uint, categoryNames []string) (*models.Image, error) {
	// 유저별 디렉토리 생성
	saveDir := fmt.Sprintf("./uploads/%d", userID)
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("서버에 이미지 디렉토리 만드는데 실패했습니다: %v", err)
		}
	}

	// 저장할 파일 경로 설정
	filePath := filepath.Join(saveDir, fileName)

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("form file을 가져오는데 실패했습니다: %v", err)
	}

	// 이미지 파일 저장
	if err := ctx.SaveUploadedFile(fileHeader, filePath); err != nil {
		return nil, fmt.Errorf("이미지를 저장하는데 실패했습니다: %v", err)
	}

	// 카테고리 검색 및 추가
	categories, err := s.categoryRepo.GetCategoriesByName(categoryNames)
	if err != nil {
		return nil, fmt.Errorf("카테고리를 가져오는데 실패했습니다: %v", err)
	}

	// 이미지 메타데이터 생성 및 저장
	image := models.Image{
		FileName:    fileName,
		FilePath:    filePath,
		UploadDate:  time.Now(),
		Description: description,
		UserID:      userID,
	}

	if err := s.imageRepo.CreateImage(&image); err != nil {
		return nil, fmt.Errorf("이미지 메타데이터를 저장하는데 실패했습니다: %v", err)
	}

	// 카테고리 매핑을 위한 image_categories 테이블 업데이트
	for _, category := range categories {
		if err := s.imageCategoryRepo.AddImageCategory(image.ID, category.ID); err != nil {
			return nil, fmt.Errorf("이미지-카테고리 매핑 업데이트를 실패했습니다: %v", err)
		}
	}

	return &image, nil
}
