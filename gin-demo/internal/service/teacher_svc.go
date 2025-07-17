package service

import (
	"context"
	"gin-demo/internal/model"
	"gin-demo/internal/repository"
)

type TeacherService struct {
	teacherRepo *repository.TeacherRepository
}

func NewTeacherService(teacherRepo *repository.TeacherRepository) *TeacherService {
	return &TeacherService{teacherRepo: teacherRepo}
}

func (s *TeacherService) TeacherRegister(ctx context.Context, req *TeacherRegisterRequest) (*model.Teacher, error) {
	teacher := &model.Teacher{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := s.teacherRepo.CreateTeacher(ctx, teacher); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *TeacherService) GetTeacherProfile(ctx context.Context, id string) (*model.Teacher, error) {
	return s.teacherRepo.GetTeacherByID(ctx, id)
}

type TeacherRegisterRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
