package service

import (
	"errors"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
)

type CourseService struct {
	courseRepo     *repository.CourseRepository
	enrollmentRepo *repository.EnrollmentRepository
}

func NewCourseService(courseRepo *repository.CourseRepository, enrollmentRepo *repository.EnrollmentRepository) *CourseService {
	return &CourseService{
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (s *CourseService) CreateCourse(course *model.Course) error {
	return s.courseRepo.Create(course)
}

func (s *CourseService) GetAllCourses() ([]model.Course, error) {
	return s.courseRepo.FindAll()
}

func (s *CourseService) GetCourseByID(id int64) (*model.Course, error) {
	return s.courseRepo.FindByID(id)
}

func (s *CourseService) Enroll(studentID, courseID int64) error {
	// 检查课程容量
	count, err := s.enrollmentRepo.CountByCourse(courseID)
	if err != nil {
		return err
	}

	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return err
	}

	if count >= int64(course.StudentMaxNum) {
		return errors.New("course is full")
	}

	return s.enrollmentRepo.Create(&model.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
	})
}

func (s *CourseService) DeleteCourse(courseID, teacherID int64) error {
	// 先检查是否有学生选课
	count, err := s.enrollmentRepo.CountByCourse(courseID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete course with enrolled students")
	}

	return s.courseRepo.DeleteByTeacher(courseID, teacherID)
}

func (s *CourseService) GetCoursesByTeacher(teacherID int64) ([]model.Course, error) {
	return s.courseRepo.FindByTeacher(teacherID)
}
