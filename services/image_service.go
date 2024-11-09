package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/repositories"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"time"
)

type ImageService interface {
	UploadImage(ctx *gin.Context, fileName, description string, userID uint, categoryNames []string) (*models.Image, error)
	GetThumbnail(imageID uint) (string, error)
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

	// 썸네일 생성
	thumbPath, err := s.createThumbnail(filePath, saveDir, fileName)
	if err != nil {
		return nil, err
	}

	// 카테고리 검색 및 추가
	categories, err := s.categoryRepo.GetCategoriesByName(categoryNames)
	if err != nil {
		return nil, fmt.Errorf("카테고리를 가져오는데 실패했습니다: %v", err)
	}

	// 이미지 메타데이터, 썸네일 경로 생성 및 저장
	uploadImage := models.Image{
		FileName:      fileName,
		FilePath:      filePath,
		ThumbnailPath: thumbPath,
		UploadDate:    time.Now(),
		Description:   description,
		UserID:        userID,
	}

	if err := s.imageRepo.CreateImage(&uploadImage); err != nil {
		return nil, fmt.Errorf("이미지 메타데이터를 저장하는데 실패했습니다: %v", err)
	}

	// 카테고리 매핑을 위한 image_categories 테이블 업데이트
	for _, category := range categories {
		if err := s.imageCategoryRepo.AddImageCategory(uploadImage.ID, category.ID); err != nil {
			return nil, fmt.Errorf("이미지-카테고리 매핑 업데이트를 실패했습니다: %v", err)
		}
	}

	return &uploadImage, nil
}

// 썸네일 생성 로직
func (s *imageService) createThumbnail(filePath, saveDir, fileName string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("썸네일 생성을 위해 이미지 파일을 여는데 실패했습니다: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("썸네일 생성을 위해 이미지 파일을 닫는데 실패했습니다: %v\n", err)
		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("썸네일 생성을 위해 이미지 디코딩이 실패했습니다: %v", err)
	}

	// 리사이즈 (예: 150x150 픽셀)
	thumbnail := resize.Thumbnail(150, 150, img, resize.Lanczos3)

	// 썸네일 저장 경로
	thumbFileName := "thumb_" + fileName
	thumbPath := filepath.Join(saveDir, thumbFileName)

	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return "", fmt.Errorf("썸네일 이미지 생성에 실패했습니다: %v", err)
	}
	defer func(thumbFile *os.File) {
		err := thumbFile.Close()
		if err != nil {
			fmt.Printf("썸네일 이미지 파일을 닫는데 실패했습니다: %v\n", err)
		}
	}(thumbFile)

	if err := jpeg.Encode(thumbFile, thumbnail, nil); err != nil {
		return "", fmt.Errorf("썸네일 저장에 실패했습니다: %v", err)
	}

	return thumbPath, nil
}

// GetThumbnail - 썸네일 경로 반환
func (s *imageService) GetThumbnail(imageID uint) (string, error) {
	uploadedImage, err := s.imageRepo.GetImageByID(imageID)
	if err != nil {
		return "", fmt.Errorf("썸네일 - 이미지를 가져오는데 실패했습니다: %v", err)
	}
	return uploadedImage.ThumbnailPath, nil
}
