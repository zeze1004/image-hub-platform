package models

import (
	"gorm.io/gorm"
	"time"
)

type Image struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	FileName      string    `gorm:"not null"`       // 원본 파일명
	FilePath      string    `gorm:"not null"`       // 서버 저장 경로
	ThumbnailPath string    `gorm:"not null"`       // 썸네일 경로
	UploadDate    time.Time `gorm:"autoCreateTime"` // 업로드된 날짜
	Description   string    // 설명
	UserID        uint      // 업로드한 사용자 ID
	gorm.Model
}
