package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/services"
	"net/http"
	"strconv"
)

type ImageController struct {
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService) *ImageController {
	return &ImageController{imageService: imageService}
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	// JWT에서 현재 요청자의 userID, Role 가져오기
	currentUserID := ctx.GetUint("userID")
	role := ctx.GetString("role")

	var targetUserID uint
	// URL 경로 파라미터로 받은 userID가 있으면, 해당 유저의 이미지를 업로드할 수 있는 ADMIN 권한이 있는지 확인
	if userIDParam := ctx.Param("userID"); userIDParam != "" {
		id, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 userID 파라미터입니다"})
			return
		}
		targetUserID = uint(id)

		if role != "ADMIN" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "어드민 권한이 없습니다"})
			return
		}
	} else {
		targetUserID = currentUserID
	}

	// 이미지 파일 및 설명 가져오기
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "이미지 파일이 누락됐습니다"})
		return
	}
	description := ctx.PostForm("description")

	// 카테고리 이름들 가져오기
	categoryNames := ctx.PostFormArray("categories")

	// 이미지 업로드
	image, err := c.imageService.UploadImage(ctx, file.Filename, description, targetUserID, categoryNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "이미지 업로드가 성공했습니다", "image": image})
}

// GetThumbnail - 이미지 ID를 받아 썸네일 이미지 파일을 반환
func (c *ImageController) GetThumbnail(ctx *gin.Context) {
	imageIDParam := ctx.Param("imageID")
	imageID, err := strconv.ParseUint(imageIDParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "imageID parameter가 잘못됐습니다"})
		return
	}

	thumbnailPath, err := c.imageService.GetThumbnail(uint(imageID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.File(thumbnailPath) // 썸네일 이미지 파일 반환
}
