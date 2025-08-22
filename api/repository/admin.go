package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin *model.Admin) error
	DeleteAdmin(jobNo string, password string) error
	UpdateAdmin(admin *model.Admin) error
	GetAdminByID(id int64) (*model.Admin, error)
	GetAdminByJobNo(jobNo string) (*model.Admin, error)
	ListAdmins() ([]model.Admin, error)
}

type GormAdminRepository struct {
	db *gorm.DB
}

func NewGormAdminRepository(db *gorm.DB) *GormAdminRepository {
	return &GormAdminRepository{db: db}
}

func (r *GormAdminRepository) CreateAdmin(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *GormAdminRepository) DeleteAdmin(jobNo string, password string) error {
	return r.db.Where("job_no = ? AND password = ?", jobNo, password).Delete(&model.Admin{}).Error
}

func (r *GormAdminRepository) UpdateAdmin(admin *model.Admin) error {
	return r.db.Model(&model.Admin{}).Where("id = ?", admin.ID).Updates(admin).Error
}

func (r *GormAdminRepository) GetAdminByID(id int64) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *GormAdminRepository) GetAdminByJobNo(jobNo string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("job_no = ?", jobNo).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *GormAdminRepository) ListAdmins() ([]model.Admin, error) {
	var admins []model.Admin
	err := r.db.Find(&admins).Error
	return admins, err
}
