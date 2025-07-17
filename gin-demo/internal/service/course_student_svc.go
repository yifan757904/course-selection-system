package service

import (
	"context"
	"gin-demo/internal/model"
	"gin-demo/internal/repository"
)

type Course_StudentService struct {
	course_studentRepo *repository.Course_StudentRepository
}

func NewCourse_StudentService(course_studentRepo *repository.Course_StudentRepository) *Course_StudentService {
	return &Course_StudentService{course_studentRepo: course_studentRepo}
}

func (s *Course_StudentService) Course_StudentRegister(ctx context.Context, req *Course_StudentRegisterRequest) (*model.Course_Student, error) {
	course_student := &model.Course_Student{
		CourseID:  req.CourseID,
		StudentID: req.StudentID,
	}

	if err := s.course_studentRepo.CreateCourse_Student(ctx, course_student); err != nil {
		return nil, err
	}

	return course_student, nil
}

func (s *Course_StudentService) GetCourse_StudentProfile(ctx context.Context, courseid string) (int, error) {
	return s.course_studentRepo.GetNumByCourseID(ctx, courseid)
}

type Course_StudentRegisterRequest struct {
	CourseID  string `json:"id" binding:"required"`
	StudentID string `json:"name" binding:"required"`
}
