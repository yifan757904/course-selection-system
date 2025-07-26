package service

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(name, idCard, rule string) (*model.User, error) {

	user := &model.User{
		Name:   name,
		IDCard: idCard,
		Rule:   rule,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func GenerateJWT(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		// 可以添加其他声明
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Login(idCard string) (string, string, error) {
	user, err := s.userRepo.FindByIDCard(idCard)
	rule := user.Rule
	if err != nil {
		return rule, "", errors.New("用户不存在")
	}

	// 确保调用了token生成
	token, err := GenerateJWT(user.ID)
	if err != nil {
		return rule, "", errors.New("token生成失败")
	}

	return rule, token, nil
}
