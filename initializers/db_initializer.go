package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/image_hub?charset=utf8&parseTime=True&loc=Local" // TODO: 환경변수로부터 DSN 정보를 가져오도록 수정
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB 연결에 실패했습니다: %v", err)
	}
	return db
}
