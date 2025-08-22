package service

import (
	"errors"
	"testing"
	"time"

	"github.com/liuyifan1996/course-selection-system/api/model"
)

type mockEnrollmentRepo struct {
	getEnrollmentFunc    func(studentID, courseID int64) (*model.Enrollment, error)
	countEnrollmentsFunc func(courseID int64) (int64, error)
}

// 实现 repository.EnrollmentRepo 接口
func (m *mockEnrollmentRepo) GetStudentByIDCard(idCard string) (*model.User, error) { return nil, nil }
func (m *mockEnrollmentRepo) GetCourseByID(courseID int) (*model.Course, error)     { return nil, nil }
func (m *mockEnrollmentRepo) GetEnrollment(studentID, courseID int64) (*model.Enrollment, error) {
	return m.getEnrollmentFunc(studentID, courseID)
}
func (m *mockEnrollmentRepo) CreateEnrollment(enrollment *model.Enrollment) error { return nil }
func (m *mockEnrollmentRepo) DeleteEnrollment(enrollment *model.Enrollment) error { return nil }
func (m *mockEnrollmentRepo) CountEnrollmentsByCourse(courseID int64) (int64, error) {
	return m.countEnrollmentsFunc(courseID)
}
func (m *mockEnrollmentRepo) GetStudentEnrollments(studentID int64) ([]model.Enrollment, error) {
	return nil, nil
}
func (m *mockEnrollmentRepo) GetStudentsByCourseID(courseID int64) ([]model.User, error) {
	return nil, nil
}
func (m *mockEnrollmentRepo) GetStudentCourses(enrollmentIDs []int64, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}

func TestCheckEnrollValid(t *testing.T) {
	svc := &EnrollmentService{}
	course := &model.Course{
		StartDate:     time.Now().Add(24 * time.Hour),
		StudentMaxNum: 2,
	}
	// 正常情况
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return nil, errors.New("not found")
		},
		countEnrollmentsFunc: func(courseID int64) (int64, error) {
			return 1, nil
		},
	}
	if err := svc.checkEnrollValid(1, 1, course); err != nil {
		t.Errorf("should pass, got %v", err)
	}
	// 已选过该课程
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return &model.Enrollment{}, nil
		},
		countEnrollmentsFunc: func(courseID int64) (int64, error) {
			return 1, nil
		},
	}
	if err := svc.checkEnrollValid(1, 1, course); err == nil || err.Error() != "已选过该课程" {
		t.Errorf("should fail for already enrolled")
	}
	// 课程已开始
	course.StartDate = time.Now().Add(-24 * time.Hour)
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return nil, errors.New("not found")
		},
		countEnrollmentsFunc: func(courseID int64) (int64, error) {
			return 1, nil
		},
	}
	if err := svc.checkEnrollValid(1, 1, course); err == nil || err.Error() != "课程已开始，不能选课" {
		t.Errorf("should fail for started course")
	}
	// 课程人数已满
	course.StartDate = time.Now().Add(24 * time.Hour)
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return nil, errors.New("not found")
		},
		countEnrollmentsFunc: func(courseID int64) (int64, error) {
			return 2, nil
		},
	}
	if err := svc.checkEnrollValid(1, 1, course); err == nil || err.Error() != "课程人数已满" {
		t.Errorf("should fail for full course")
	}
	// 获取人数失败
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return nil, errors.New("not found")
		},
		countEnrollmentsFunc: func(courseID int64) (int64, error) {
			return 0, errors.New("db error")
		},
	}
	if err := svc.checkEnrollValid(1, 1, course); err == nil || err.Error() != "无法获取课程人数" {
		t.Errorf("should fail for count error")
	}
}

func TestCheckUnenrollValid(t *testing.T) {
	svc := &EnrollmentService{}
	course := &model.Course{
		StartDate: time.Now().Add(24 * time.Hour),
	}
	// 正常情况
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return &model.Enrollment{}, nil
		},
	}
	if err := svc.checkUnenrollValid(1, 1, course); err != nil {
		t.Errorf("should pass, got %v", err)
	}
	// 未选课
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return nil, errors.New("not found")
		},
	}
	if err := svc.checkUnenrollValid(1, 1, course); err == nil || err.Error() != "未选择该课程" {
		t.Errorf("should fail for not enrolled")
	}
	// 课程已开始
	course.StartDate = time.Now().Add(-24 * time.Hour)
	svc.repo = &mockEnrollmentRepo{
		getEnrollmentFunc: func(studentID, courseID int64) (*model.Enrollment, error) {
			return &model.Enrollment{}, nil
		},
	}
	if err := svc.checkUnenrollValid(1, 1, course); err == nil || err.Error() != "课程已开始，不能退选" {
		t.Errorf("should fail for started course")
	}
}
