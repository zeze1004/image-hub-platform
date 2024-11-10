package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/services"
	"github.com/zeze1004/image-hub-platform/utils"
	"net/http"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// GetCategoriesByImageID 특정 이미지의 카테고리 조회
func (c *CategoryController) GetCategoriesByImageID(ctx *gin.Context) {
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	categories, err := c.categoryService.GetCategoriesByImageIDAndUserID(imageID, ctx.GetUint("userID"), c.isAdmin(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categories)
	return
}

// AddCategoryToImage 이미지에 카테고리 추가
func (c *CategoryController) AddCategoryToImage(ctx *gin.Context) {
	categoryIDParam := ctx.Param("categoryID")
	categoryID, _ := c.parseAndValidateID(categoryIDParam)
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	if err := c.categoryService.AddCategoryToImageByImageIDAndCategoryID(imageID, categoryID, ctx.GetUint("userID"), c.isAdmin(ctx)); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "이미지에 카테고리가 추가됐습니다"})
}

// RemoveCategoryFromImage 이미지에서 카테고리 삭제
func (c *CategoryController) RemoveCategoryFromImage(ctx *gin.Context) {
	categoryIDParam := ctx.Param("categoryID")
	categoryID, _ := c.parseAndValidateID(categoryIDParam)
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	if err := c.categoryService.RemoveCategoryFromImageByImageIDAndCategoryID(imageID, categoryID, ctx.GetUint("userID"), c.isAdmin(ctx)); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "이미지에서 카테고리를 삭제했습니다"})
}

// isAdmin 관리자 권한인지 확인
func (c *CategoryController) isAdmin(ctx *gin.Context) bool {
	return utils.IsAdmin(ctx)
}

// parseAndValidateID userID 파라미터 파싱 및 유효성 검사
func (c *CategoryController) parseAndValidateID(paramID string) (uint, error) {
	return utils.ParseAndValidateID(paramID)
}
