package service

import (
	"fmt"
	"time"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"github.com/liuyifan1996/course-selection-system/api/repository"
)

type EnrollmentService struct {
	repo *repository.EnrollmentRepository
}

func NewEnrollmentService(repo *repository.EnrollmentRepository) *EnrollmentService {
	return &EnrollmentService{repo: repo}
}

func (s *EnrollmentService) Enroll(studentIDCard string, courseID int) error {
	// 检查学生是否存在
	student, err := s.repo.GetStudentByIDCard(studentIDCard)
	if err != nil {
		return fmt.Errorf("学生不存在")
	}

	// 检查是否已选课
	existing, err := s.repo.GetEnrollment(student.ID, int64(courseID))
	if err == nil && existing != nil {
		return fmt.Errorf("已选过该课程")
	}

	// 检查课程是否存在
	course, err := s.repo.GetCourseByID(courseID)
	if err != nil {
		return fmt.Errorf("课程不存在")
	}

	// 检查课程是否已开始
	if course.StartDate.Before(time.Now()) {
		return fmt.Errorf("课程已开始，不能选课")
	}

	// 检查课程是否已满
	count, err := s.repo.CountEnrollmentsByCourse(int64(courseID))
	if err != nil {
		return fmt.Errorf("无法获取课程人数")
	}
	if count >= int64(course.StudentMaxNum) {
		return fmt.Errorf("课程人数已满")
	}

	enrollment := &model.Enrollment{
		StudentID: int64(student.ID),
		CourseID:  int64(courseID),
	}

	if err := s.repo.CreateEnrollment(enrollment); err != nil {
		return fmt.Errorf("选课失败: %v", err)
	}

	return nil
}

func (s *EnrollmentService) GetStudentCourses(studentIDCard string, page, pageSize int, sortBy, sortOrder string, fields []string) (*model.PaginatedResponse[map[string]interface{}], error) {
	// 检查学生是否存在
	student, err := s.repo.GetStudentByIDCard(studentIDCard)
	if err != nil {
		return nil, fmt.Errorf("学生不存在")
	}

	// 查找学生的选课记录
	enrollments, err := s.repo.GetStudentEnrollments(student.ID)
	if err != nil {
		return nil, fmt.Errorf("获取选课记录失败")
	}

	// 保存所选课程id
	var courseIDs []int64
	for _, e := range enrollments {
		courseIDs = append(courseIDs, e.CourseID)
	}

	// 设置分页
	pagination := model.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	// 获取课程列表
	courses, total, err := s.repo.GetStudentCourses(courseIDs, pagination, sortBy, sortOrder, fields)
	if err != nil {
		return nil, fmt.Errorf("查询课程失败")
	}

	// 构建响应
	response := &model.PaginatedResponse[map[string]interface{}]{
		Data:       courses,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)),
	}

	return response, nil
}

func (s *EnrollmentService) DeleteEnrollment(studentIDCard string, courseID int) error {
	// 检查学生是否存在
	student, err := s.repo.GetStudentByIDCard(studentIDCard)
	if err != nil {
		return fmt.Errorf("学生不存在")
	}

	// 检查是否已选课
	existing, err := s.repo.GetEnrollment(student.ID, int64(courseID))
	if err != nil {
		return fmt.Errorf("未选择该课程")
	}

	// 检查课程是否存在
	course, err := s.repo.GetCourseByID(courseID)
	if err != nil {
		return fmt.Errorf("课程不存在")
	}

	// 检查课程是否已开始
	if course.StartDate.Before(time.Now()) {
		return fmt.Errorf("课程已开始，不能退选")
	}

	// 删除选课
	if err := s.repo.DeleteEnrollment(existing); err != nil {
		return fmt.Errorf("退选失败: %v", err)
	}

	return nil
}
