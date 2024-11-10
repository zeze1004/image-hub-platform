package test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeze1004/image-hub-platform/mocks"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/services"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// 이미지 업로드 성공 테스트 - 해당 테스트를 실패하면 test/uploads 디렉토리에 테스트 이미지 파일이 생성됨
func TestUploadImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := mocks.NewMockImageRepository(ctrl)
	mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)
	mockImageCategoryRepo := mocks.NewMockImageCategoryRepository(ctrl)

	mockImageRepo.EXPECT().CreateImage(gomock.Any()).DoAndReturn(func(image *models.Image) error {
		image.ID = 1
		return nil
	})

	mockCategoryRepo.EXPECT().GetCategoriesByName(gomock.Any()).Return([]models.Category{{ID: 1, Name: "TEST_CATEGORY"}}, nil)

	mockImageCategoryRepo.EXPECT().AddImageCategory(gomock.Any(), gomock.Any()).Return(nil)

	imageService := services.NewImageService(mockImageRepo, mockCategoryRepo, mockImageCategoryRepo)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)

	saveDir := "./uploads/1"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Errorf("테스트 파일을 삭제하는데 실패했습니다: %v", err)
		}
	}(saveDir)

	fileName := "test.jpg"
	description := "test image"
	userID := uint(1)
	categoryNames := []string{"TEST_CATEGORY"}

	// 테스트 이미지 생성
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, nil)
	if err != nil {
		t.Fatalf("이미지 인코딩에 실패했습니다: %v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		t.Fatalf("form file을 만드는데 실패했습니다: %v", err)
	}
	_, err = part.Write(imgBuf.Bytes())
	if err != nil {
		t.Fatalf("form file에 작성하는데 실패했습니다: %v", err)
	}
	writer.Close()

	ctx.Request = &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(&buf),
	}
	ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())

	uploadedImage, err := imageService.UploadImage(ctx, fileName, description, userID, categoryNames)

	assert.NoError(t, err)
	assert.Equal(t, fileName, uploadedImage.FileName)
	assert.Equal(t, filepath.Join(saveDir, fileName), uploadedImage.FilePath)
	assert.NotEmpty(t, uploadedImage.ThumbnailPath)
	assert.Equal(t, description, uploadedImage.Description)
	assert.Equal(t, userID, uploadedImage.UserID)
	assert.WithinDuration(t, time.Now(), uploadedImage.UploadDate, time.Second)
}

// 카테고리가 달라 이미지 업로드 실패 케이스 테스트
func TestUploadImageFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockImageRepo := mocks.NewMockImageRepository(ctrl)
	mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)
	mockImageCategoryRepo := mocks.NewMockImageCategoryRepository(ctrl)

	// 카테고리가 달라서 이미지 생성 실패
	mockCategoryRepo.EXPECT().GetCategoriesByName(gomock.Any()).Return([]models.Category{{ID: 1, Name: "TEST_CATEGORY"}}, nil)
	mockImageRepo.EXPECT().CreateImage(gomock.Any()).Return(fmt.Errorf("이미지 생성에 실패했습니다"))

	imageService := services.NewImageService(mockImageRepo, mockCategoryRepo, mockImageCategoryRepo)

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)

	saveDir := "./test_uploads/1"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Errorf("테스트 파일을 삭제하는데 실패했습니다: %v", err)
		}
	}(saveDir)

	fileName := "test.jpg"
	description := "test image"
	userID := uint(1)
	categoryNames := []string{"WRONG_CATEGORY"} // 잘못된 카테고리

	// 테스트 이미지 생성
	img := image.NewRGBA(image.Rect(0, 0, 100, 100)) // 100x100 픽셀 이미지 생성
	var imgBuf bytes.Buffer
	err := jpeg.Encode(&imgBuf, img, nil)
	if err != nil {
		t.Fatalf("이미지 인코딩에 실패했습니다: %v", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		t.Fatalf("form file을 만드는데 실패했습니다: %v", err)
	}
	_, err = part.Write(imgBuf.Bytes())
	if err != nil {
		t.Fatalf("form file에 작성하는데 실패했습니다: %v", err)
	}
	writer.Close()

	ctx.Request = &http.Request{
		Header: make(http.Header),
		Body:   io.NopCloser(&buf),
	}
	ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())

	uploadedImage, err := imageService.UploadImage(ctx, fileName, description, userID, categoryNames)

	assert.Error(t, err)
	assert.Nil(t, uploadedImage)
	assert.Contains(t, err.Error(), "이미지 생성에 실패했습니다")
}

// 모든 이미지 파일 삭제를 병렬로 처리할 떄와 순차 삭제 성능 비교
func BenchmarkDeleteAllImages(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockImageRepo := mocks.NewMockImageRepository(ctrl)
	mockCategoryRepo := mocks.NewMockCategoryRepository(ctrl)
	mockImageCategoryRepo := mocks.NewMockImageCategoryRepository(ctrl)

	imageService := services.NewImageService(mockImageRepo, mockCategoryRepo, mockImageCategoryRepo)

	gin.SetMode(gin.TestMode)

	userID := uint(1)
	numFiles := 1000
	images := make([]models.Image, numFiles)
	// 테스트 이미지 파일 이름 생성
	for i := range images {
		images[i] = models.Image{
			ID:            userID,
			FilePath:      "./uploads/" + strconv.Itoa(i+1) + ".jpg",
			ThumbnailPath: "./uploads/thumb_" + strconv.Itoa(i+1) + ".jpg",
		}
	}

	// 필요 의존성 메소드 MOCK 처리
	mockImageRepo.EXPECT().GetImagesByUserID(userID).Return(images, nil).AnyTimes()
	mockImageRepo.EXPECT().DeleteImagesByUserID(userID).Return(nil).AnyTimes()

	for _, bm := range []struct {
		name string
		fn   func() error
	}{
		{
			name: "ParallelDelete",
			// 병렬로 모든 이미지 삭제하는 DeleteAllImagesByUserID 호출
			fn: func() error {
				return imageService.DeleteAllImagesByUserID(userID)
			},
		},
		{
			name: "SequentialDelete",
			// 모든 파일을 순차적으로 삭제하는 반복문 호출
			fn: func() error {
				for _, image := range images {
					if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
						return err
					}
					if err := os.Remove(image.ThumbnailPath); err != nil && !os.IsNotExist(err) {
						return err
					}
				}
				return nil
			},
		},
	} {
		b.Run(bm.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ { // 실행 속도 평균을 내기 위해 반복 실행
				if err := createTestFiles(images); err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				if err := bm.fn(); err != nil {
					b.Fatal(err)
				}
				b.StopTimer()
			}
		})
	}
}

// 테스트 이미지 파일 생성
func createTestFiles(images []models.Image) error {
	for _, image := range images {
		if err := os.MkdirAll(filepath.Dir(image.FilePath), os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(image.FilePath)
		if err != nil {
			return err
		}
		_, err = file.Write([]byte(image.FileName))
		if err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}
