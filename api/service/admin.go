package service

import (
	"errors"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
)

type AdminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

type CreateAdminInput struct {
	JobNo    string
	Password string
}

type UpdateAdminInput struct {
	ID       int64
	JobNo    *string
	Password *string
}

func (s *AdminService) CreateAdmin(input CreateAdminInput) (*model.Admin, error) {
	admin := &model.Admin{
		JobNo:    input.JobNo,
		Password: input.Password,
	}
	if err := s.repo.CreateAdmin(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

type DeleteAdminInput struct {
	JobNo    string
	Password string
}

func (s *AdminService) DeleteAdmin(input DeleteAdminInput) error {
	return s.repo.DeleteAdmin(input.JobNo, input.Password)
}

func (s *AdminService) UpdateAdmin(input UpdateAdminInput) (*model.Admin, error) {
	admin, err := s.repo.GetAdminByID(input.ID)
	if err != nil {
		return nil, err
	}
	if input.JobNo != nil {
		admin.JobNo = *input.JobNo
	}
	if input.Password != nil {
		admin.Password = *input.Password
	}
	if err := s.repo.UpdateAdmin(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *AdminService) GetAdminByID(id int64) (*model.Admin, error) {
	return s.repo.GetAdminByID(id)
}

func (s *AdminService) GetAdminByJobNo(jobNo string) (*model.Admin, error) {
	return s.repo.GetAdminByJobNo(jobNo)
}

func (s *AdminService) ListAdmins() ([]model.Admin, error) {
	return s.repo.ListAdmins()
}

var ErrAdminInvalidCredentials = errors.New("工号或密码错误")

type AdminLoginInput struct {
	JobNo    string
	Password string
}

func (s *AdminService) Login(input AdminLoginInput) (*model.Admin, error) {
	admin, err := s.repo.GetAdminByJobNo(input.JobNo)
	if err != nil || admin.Password != input.Password {
		return nil, ErrAdminInvalidCredentials
	}
	return admin, nil
}
