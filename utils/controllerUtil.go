package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// IsAdmin 관리자 권한인지 확인
func IsAdmin(ctx *gin.Context) bool {
	role := ctx.GetString("role")
	return role == "ADMIN"
}

// ParseAndValidateID userID 파라미터 파싱 및 유효성 검사
func ParseAndValidateID(paramID string) (uint, error) {
	if paramID == "" {
		return 0, fmt.Errorf("파라미터가 비어 있습니다")
	}
	id, err := strconv.ParseUint(paramID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("잘못된 파라미터입니다")
	}
	return uint(id), nil
}
