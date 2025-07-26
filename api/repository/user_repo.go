package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByIDCard(idCard string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id_card = ?", idCard).First(&user).Error
	return &user, err
}
