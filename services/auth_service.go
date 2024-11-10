package services

import (
	"errors"
	"github.com/zeze1004/image-hub-platform/models"
	"github.com/zeze1004/image-hub-platform/repositories"
	"github.com/zeze1004/image-hub-platform/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(user *models.User) error
	Login(email, password string) (string, error) // TODO: JWT 토큰 리턴 수정
}

// authService 구조체
type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepo}
}

// SignUp - 새로운 사용자 등록
func (s *authService) SignUp(user *models.User) error {
	// TODO: 이메일 중복시 클라이언트단에 알림
	_, err := s.userRepository.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("이미 등록된 이메일입니다")
	}

	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 사용자 저장
	return s.userRepository.CreateUser(user)
}

// Login - 이메일과 비밀번호로 로그인 처리
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("가입되지 않은 메일입니다")
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("비밀번호가 일치하지 않습니다")
	}

	// JWT 토큰 생성
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
