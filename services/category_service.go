package services

import (
	"fmt"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/repositories"
	"github.com/zeze1004/image-hub-platform/utils"
)

type categoryService struct {
	imageRepo         repositories.ImageRepository
	categoryRepo      repositories.CategoryRepository
	imageCategoryRepo repositories.ImageCategoryRepository
}

func NewCategoryService(imageRepo repositories.ImageRepository, categoryRepo repositories.CategoryRepository, imageCategoryRepo repositories.ImageCategoryRepository) CategoryService {
	return &categoryService{imageRepo: imageRepo, categoryRepo: categoryRepo, imageCategoryRepo: imageCategoryRepo}
}

// GetCategoriesByImageIDAndUserID 특정 이미지의 카테고리 조회
func (s *categoryService) GetCategoriesByImageIDAndUserID(imageID, userID uint, isAdmin bool) ([]models.Category, error) {
	if !isAdmin {
		if err := s.validateImageOwnership(imageID, userID); err != nil {
			return nil, err
		}
	}

	return s.categoryRepo.GetCategoriesByImageID(imageID)
}

// AddCategoryToImageByImageIDAndCategoryID 이미지에 카테고리 추가
func (s *categoryService) AddCategoryToImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error {
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
func (s *categoryService) RemoveCategoryFromImageByImageIDAndCategoryID(imageID, categoryID, userID uint, isAdmin bool) error {
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

func (s *categoryService) isDuplicateCategory(imageID uint, categoryID uint) error {
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

func (s *categoryService) isValidCategoryID(categoryID uint) error {
	if categoryID > 5 || categoryID < 1 {
		return fmt.Errorf("잘못된 카테고리 ID입니다")
	}
	return nil
}

func (s *categoryService) validateImageOwnership(imageID, userID uint) error {
	return utils.ValidateImageOwnership(s.imageRepo, imageID, userID)
}
