package initializers

import (
	"github.com/zeze1004/image-hub-platform/controllers"
	"github.com/zeze1004/image-hub-platform/repositories"
	"github.com/zeze1004/image-hub-platform/services"
	"gorm.io/gorm"
)

func InitUserModule(db *gorm.DB) *controllers.AuthController {
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)
	return authController
}
