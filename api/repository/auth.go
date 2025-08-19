package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByIDCard(idCard string) (*model.User, error)
	CreateUser(user *model.User) error
}

type GormAuthRepository struct {
	db *gorm.DB
}

func NewGormAuthRepository(db *gorm.DB) *GormAuthRepository {
	return &GormAuthRepository{db: db}
}

func (r *GormAuthRepository) FindByIDCard(idCard string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id_card = ?", idCard).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormAuthRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
