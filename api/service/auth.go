package service

import (
	"errors"
	"regexp"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
	"github.com/liuyifan1996/course-selection-system/pkg"
)

var (
	ErrUserAlreadyExists  = errors.New("用户已存在")
	ErrInvalidPassword    = errors.New("密码必须包含至少一个大写字母、一个小写字母和一个数字")
	ErrInvalidCredentials = errors.New("用户不存在或密码错误")
)

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

type RegisterInput struct {
	IDCard   string `json:"id_card"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"` // student or teacher
}

func (s *AuthService) Register(input RegisterInput) (*model.User, error) {
	// 检查用户是否已存在
	if _, err := s.repo.FindByIDCard(input.IDCard); err == nil {
		return nil, ErrUserAlreadyExists
	}

	// 密码复杂度验证
	if !isValidPassword(input.Password) {
		return nil, ErrInvalidPassword
	}

	user := &model.User{
		IDCard:   input.IDCard,
		Name:     input.Name,
		Password: input.Password,
		Role:     input.Role,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	IDCard   string `json:"id_card"`
	Password string `json:"password"`
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	user, err := s.repo.FindByIDCard(input.IDCard)
	if err != nil || user.Password != input.Password {
		return "", ErrInvalidCredentials
	}

	token, err := pkg.GenerateToken(user.IDCard, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func isValidPassword(password string) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(password)
}
