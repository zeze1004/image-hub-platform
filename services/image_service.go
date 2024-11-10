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
	"sync"
	"time"
)

type ImageService interface {
	UploadImage(ctx *gin.Context, fileName, description string, userID uint, categoryNames []string) (*models.Image, error)
	GetThumbnail(imageID uint) (string, error)
	GetAllImages() ([]models.Image, error)
	GetImagesByUserID(userID uint) ([]models.Image, error)
	GetImageByID(imageID uint, userID uint, isAdmin bool) (*models.Image, error)
	DeleteImageByID(imageID uint, userID uint, isAdmin bool) error
	DeleteAllImagesByUserID(userID uint) error
	GetImagesByCategoryIDAndUserID(categoryID, userID uint, isAdmin bool) ([]models.Image, error)  // (유저) 특정 카테고리를 갖는 사용자의 이미지 조회
	GetCategoriesByImageIDAndUserID(imageID, userID uint, isAdmin bool) ([]models.Category, error) // (유저) 특정 이미지의 카테고리 조회
	AddCategoryToImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error
	RemoveCategoryFromImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error
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

// GetThumbnail 썸네일 경로 반환
func (s *imageService) GetThumbnail(imageID uint) (string, error) {
	uploadedImage, err := s.imageRepo.GetImageByID(imageID)
	if err != nil {
		return "", fmt.Errorf("썸네일 - 이미지를 가져오는데 실패했습니다: %v", err)
	}
	return uploadedImage.ThumbnailPath, nil
}

// GetAllImages 모든 이미지 목록 조회
func (s *imageService) GetAllImages() ([]models.Image, error) {
	images, err := s.imageRepo.GetAllImages()
	if err != nil {
		return nil, fmt.Errorf("이미지 목록을 가져오는데 실패했습니다: %v", err)
	}
	return images, nil
}

// GetImagesByUserID - 특정 사용자의 이미지 목록 조회
func (s *imageService) GetImagesByUserID(userID uint) ([]models.Image, error) {
	images, err := s.imageRepo.GetImagesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("유저의 이미지 목록을 가져오는데 실패했습니다: %v", err)
	}
	return images, nil
}

// GetImageByID imageID로 이미지 조회
func (s *imageService) GetImageByID(imageID uint, userID uint, isAdmin bool) (*models.Image, error) {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return nil, err
		}
	}

	return s.imageRepo.GetImageByID(imageID)
}

// DeleteImageByID imageID로 개별 이미지 삭제
func (s *imageService) DeleteImageByID(imageID uint, userID uint, isAdmin bool) error {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return err
		}
	}

	return s.imageRepo.DeleteImage(imageID)
}

// DeleteAllImagesByUserID user의 모든 이미지 삭제
func (s *imageService) DeleteAllImagesByUserID(userID uint) error {
	images, err := s.imageRepo.GetImagesByUserID(userID)
	if err != nil {
		return err
	}

	// 데이터베이스에서 이미지 데이터 일괄 삭제
	if err := s.imageRepo.DeleteImagesByUserID(userID); err != nil {
		return fmt.Errorf("DB에서 이미지 삭제가 실패했습니다: %v", err)
	}

	// 파일 시스템에서 이미지 파일 및 썸네일 병렬 삭제
	errChan := make(chan error, len(images)) // 고루틴 에러를 수집할 채널
	var wg sync.WaitGroup

	for _, image := range images {
		wg.Add(1)
		go func(img models.Image) {
			defer wg.Done()
			if err := s.deleteImageFiles(&img); err != nil {
				errChan <- err // 에러 발생 시 채널에 전송
			}
		}(image)
	}

	// 모든 고루틴이 종료되면 채널 닫기
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 첫 번째 에러 반환
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

// deleteImageFiles 파일 시스템에서 이미지 및 썸네일 파일 삭제
func (s *imageService) deleteImageFiles(image *models.Image) error {
	var deleteErrors []error

	// 이미지 파일 삭제
	if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
		deleteErrors = append(deleteErrors, fmt.Errorf("서버에서 이미지 삭제를 실패했습니다: %v", err))
	}

	// 썸네일 파일 삭제
	if err := os.Remove(image.ThumbnailPath); err != nil && !os.IsNotExist(err) {
		deleteErrors = append(deleteErrors, fmt.Errorf("서버에서 썸네일 삭제를 실패했습니다: %v", err))
	}

	// 에러가 있으면 첫 번째 에러 반환
	if len(deleteErrors) > 0 {
		return deleteErrors[0]
	}
	return nil
}

// GetImagesByCategoryIDAndUserID 특정 카테고리를 가진 사용자의 이미지 조회
func (s *imageService) GetImagesByCategoryIDAndUserID(categoryID, userID uint, isAdmin bool) ([]models.Image, error) {
	if !isAdmin {
		return s.categoryRepo.GetImagesByCategoryIDAndUserID(categoryID, userID)
	} else {
		return s.categoryRepo.GetImagesByCategoryID(categoryID)
	}
}

// GetCategoriesByImageIDAndUserID 특정 이미지의 카테고리 조회
func (s *imageService) GetCategoriesByImageIDAndUserID(imageID, userID uint, isAdmin bool) ([]models.Category, error) {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return nil, err
		}
	}

	return s.categoryRepo.GetCategoriesByImageID(imageID)
}

// AddCategoryToImageByImageIDAndCategoryID 이미지에 카테고리 추가
func (s *imageService) AddCategoryToImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return err
		}
	}

	err := s.isValidCategoryID(categoryID)
	if err != nil {
		return err
	}

	// 추가할 카테고리가 이미 이미지에 있는지 검증
	err = s.isDuplicateCategory(imageID, categoryID)
	if err != nil {
		return err
	}

	return s.imageCategoryRepo.AddCategoryToImage(imageID, categoryID)
}

// RemoveCategoryFromImageByImageIDAndCategoryID 이미지에서 카테고리 제거
func (s *imageService) RemoveCategoryFromImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return err
		}
	}

	err := s.isValidCategoryID(categoryID)
	if err != nil {
		return err
	}

	err = s.isDuplicateCategory(imageID, categoryID)
	// 중복 카테고리가 없으면 에러가 반환되지 않으므로 삭제할 카테고리가 없다는 에러 반환
	if err == nil {
		return fmt.Errorf("이미지에 카테고리가 없어서 삭제할 수 없습니다") // TODO: 500 에러 리턴되는 로직 수정
	}

	return s.imageCategoryRepo.RemoveCategoryFromImage(imageID, categoryID)
}

func (s *imageService) isDuplicateCategory(imageID uint, categoryID uint) error {
	// 카테고리 중복 검증
	categories, err := s.imageCategoryRepo.GetCategoriesByImageID(imageID)
	if err != nil {
		return err
	}
	for _, category := range categories {
		if category.ID == categoryID {
			return fmt.Errorf("이미지에 이미 등록된 카테고리입니다")
		}
	}
	return nil
}

func (s *imageService) isValidCategoryID(categoryID uint) error {
	if categoryID > 5 || categoryID < 1 {
		return fmt.Errorf("잘못된 카테고리 ID입니다")
	}
	return nil
}

// validateImageOwnership 이미지 소유권 검증
func (s *imageService) validateImageOwnership(imageID, userID uint) error {
	image, err := s.imageRepo.GetImageByID(imageID)
	if err != nil {
		return fmt.Errorf("이미지를 찾을 수 없습니다: %v", err)
	}
	if image.UserID != userID {
		return fmt.Errorf("이미지에 대한 권한이 없습니다")
	}
	return nil
}
