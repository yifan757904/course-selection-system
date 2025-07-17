package service

import (
	"context"

	"github.com/liuyifan1996/internal/model"
	"github.com/liuyifan1996/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// 实际项目中应该对密码进行哈希处理
	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 注意: 生产环境必须加密
		Email:    req.Email,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, userID int64) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
