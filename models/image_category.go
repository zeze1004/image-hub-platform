package models

import "time"

type ImageCategory struct {
	ImageID    uint `gorm:"column:image_id"`
	CategoryID uint `gorm:"column:category_id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (ImageCategory) TableName() string {
	return "image_categories"
}
