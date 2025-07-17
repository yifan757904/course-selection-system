package service

import (
	"context"
	"gin-demo/internal/model"
	"gin-demo/internal/repository"
)

type CourseService struct {
	courseRepo *repository.CourseRepository
}

func NewCourseService(courseRepo *repository.CourseRepository) *CourseService {
	return &CourseService{courseRepo: courseRepo}
}

func (s *CourseService) CourseRegister(ctx context.Context, req *CourseRegisterRequest) (*model.Course, error) {
	course := &model.Course{
		ID:             req.ID,
		Name:           req.Name,
		TeacherID:      req.TeacherID,
		TeacherName:    req.TeacherName,
		Remarks:        req.Remarks,
		Student_maxnum: req.Student_maxnum,
		Time_max:       req.Time_max,
		Time_min:       req.Time_min,
	}

	if err := s.courseRepo.CreateCourse(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (s *CourseService) GetCourseProfile(ctx context.Context, id string) (*model.Course, error) {
	return s.courseRepo.GetCourseByID(ctx, id)
}

type CourseRegisterRequest struct {
	ID             string `json:"id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	TeacherID      string `json:"teacher_id" binding:"required"`
	TeacherName    string `json:"teacher_name" binding:"required"`
	Remarks        string `json:"remarks"`
	Student_maxnum int    `json:"student_maxnum" binding:"required"`
	Time_max       int    `json:"time_max"`
	Time_min       int    `json:"time_min"`
}
