package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/services"
	"net/http"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) SignUp(ctx *gin.Context) {
	var userReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.SignUp(&models.User{
		Email:    userReq.Email,
		Password: userReq.Password,
		Role:     "USER",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "회원가입이 실패됐습니다"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "인증 실패했습니다"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
