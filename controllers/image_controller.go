package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/services"
	"github.com/zeze1004/image-hub-platform/utils"
	"net/http"
)

type ImageController struct {
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService) *ImageController {
	return &ImageController{imageService: imageService}
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	var userID uint
	// ADMIN 권한 요청이라면, URL 경로 파라미터로 받은 userID가 있는지 검증
	if c.isAdmin(ctx) {
		userIDParam := ctx.Param("userID")
		userID, _ = c.parseAndValidateID(userIDParam)
	} else {
		userID = ctx.GetUint("userID")
	}

	// 이미지 파일, 설명, 카테고리 가져오기
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "이미지 파일이 누락됐습니다"})
		return
	}
	description := ctx.PostForm("description")
	categoryNames := ctx.PostFormArray("categories")

	// 이미지 업로드
	image, err := c.imageService.UploadImage(ctx, file.Filename, description, userID, categoryNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "이미지 업로드가 성공했습니다", "image": image})
}

// GetThumbnail - 이미지 ID를 받아 썸네일 이미지 파일을 반환
func (c *ImageController) GetThumbnail(ctx *gin.Context) {
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	thumbnailPath, err := c.imageService.GetThumbnail(imageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(thumbnailPath) // 썸네일 이미지 파일 반환
}

// GetImagesByUserID User가 가진 모든 이미지 목록 조회
func (c *ImageController) GetImagesByUserID(ctx *gin.Context) {
	var userID uint
	if c.isAdmin(ctx) {
		userIDParam := ctx.Param("userID")
		userID, _ = c.parseAndValidateID(userIDParam)
	} else {
		userID = ctx.GetUint("userID")
	}
	images, _ := c.imageService.GetImagesByUserID(userID)
	ctx.JSON(http.StatusOK, images)
}

// GetImageByID imageID로 특정 이미지 조회
func (c *ImageController) GetImageByID(ctx *gin.Context) {
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	image, err := c.imageService.GetImageByID(imageID, ctx.GetUint("userID"), c.isAdmin(ctx))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, image)
}

// GetAllImagesByAdmin 모든 이미지 목록 조회
func (c *ImageController) GetAllImagesByAdmin(ctx *gin.Context) {
	images, err := c.imageService.GetAllImages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, images)

}

// DeleteImage 개별 이미지 삭제
func (c *ImageController) DeleteImage(ctx *gin.Context) {
	imageIDParam := ctx.Param("imageID")
	imageID, _ := c.parseAndValidateID(imageIDParam)

	if err := c.imageService.DeleteImageByID(imageID, ctx.GetUint("userID"), c.isAdmin(ctx)); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "이미지 삭제가 성공했습니다"})
}

// DeleteAllUserImages 유저의 모든 이미지 삭제
func (c *ImageController) DeleteAllUserImages(ctx *gin.Context) {
	var userID uint
	if c.isAdmin(ctx) {
		userIDParam := ctx.Param("userID")
		userID, _ = c.parseAndValidateID(userIDParam)
	} else {
		userID = ctx.GetUint("userID")
	}

	if err := c.imageService.DeleteAllImagesByUserID(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "모든 이미지 삭제가 성공했습니다"})
}

// GetImagesByCategoryID - 특정 카테고리를 갖는 이미지 조회
func (c *ImageController) GetImagesByCategoryID(ctx *gin.Context) {
	categoryIDParam := ctx.Param("categoryID")
	categoryID, _ := c.parseAndValidateID(categoryIDParam)

	images, err := c.imageService.GetImagesByCategoryIDAndUserID(categoryID, ctx.GetUint("userID"), c.isAdmin(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, images)
}

// isAdmin 관리자 권한인지 확인
func (c *ImageController) isAdmin(ctx *gin.Context) bool {
	return utils.IsAdmin(ctx)
}

// parseAndValidateID userID 파라미터 파싱 및 유효성 검사
func (c *ImageController) parseAndValidateID(paramID string) (uint, error) {
	return utils.ParseAndValidateID(paramID)
}
