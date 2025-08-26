package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByIDCard(idCard string) (*model.User, error)
	CreateUser(user *model.User) error
	DeleteUser(idCard string, password string) error
	UpdateUser(user *model.User) error
	GetUserByID(id int64) (*model.User, error)
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

func (r *GormAuthRepository) DeleteUser(idCard string, password string) error {
	return r.db.Where("id_card = ? AND password = ?", idCard, password).Delete(&model.User{}).Error
}

func (r *GormAuthRepository) UpdateUser(user *model.User) error {
	return r.db.Model(&model.User{}).Where("id_card = ?", user.IDCard).Updates(user).Error
}

func (r *GormAuthRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
