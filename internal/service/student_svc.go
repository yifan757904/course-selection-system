package service

import (
	"context"

	"github.com/liuyifan1996/internal/model"
	"github.com/liuyifan1996/internal/repository"
)

type StudentService struct {
	studentRepo *repository.StudentRepository
}

func NewStudentService(studentRepo *repository.StudentRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo}
}

func (s *StudentService) StudentRegister(ctx context.Context, req *StudentRegisterRequest) (*model.Student, error) {
	student := &model.Student{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := s.studentRepo.CreateStudent(ctx, student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *StudentService) GetStudentProfile(ctx context.Context, id string) (*model.Student, error) {
	return s.studentRepo.GetStudentByID(ctx, id)
}

type StudentRegisterRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
