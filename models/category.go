package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	name string `gorm:"unique;not null"`
}
