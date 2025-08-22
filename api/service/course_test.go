package service

import (
	"errors"
	"testing"

	"github.com/liuyifan1996/course-selection-system/api/model"
)

type mockCourseRepo struct {
	courses            map[int64]*model.Course
	createErr          error
	updateErr          error
	getByIDErr         error
	getEnrollmentCount int64
}

func (m *mockCourseRepo) Create(course *model.Course) error { return m.createErr }
func (m *mockCourseRepo) Update(course *model.Course, data map[string]interface{}) error {
	return m.updateErr
}
func (m *mockCourseRepo) GetByID(id int64) (*model.Course, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	c, ok := m.courses[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}
func (m *mockCourseRepo) GetEnrollmentCount(courseID int64) (int64, error) {
	return m.getEnrollmentCount, nil
}
func (m *mockCourseRepo) Delete(courseID int64) error { return nil }
func (m *mockCourseRepo) GetAll(p model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}
func (m *mockCourseRepo) GetByTeacherID(teacherID string, p model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}
func (m *mockCourseRepo) GetByTeacherName(teacherName string, p model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}
func (m *mockCourseRepo) GetByCourseName(courseName string, p model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}

type mockUserRepo struct {
	teachers map[string]*model.User
}

func (m *mockUserRepo) FindByIDCard(idCard string) (*model.User, error) {
	t, ok := m.teachers[idCard]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}
func (m *mockUserRepo) CreateUser(user *model.User) error         { return nil }
func (m *mockUserRepo) DeleteUser(idCard, password string) error  { return nil }
func (m *mockUserRepo) UpdateUser(user *model.User) error         { return nil }
func (m *mockUserRepo) GetUserByID(id int64) (*model.User, error) { return nil, nil }

func TestCreateCourse_TeacherNotFound(t *testing.T) {
	courseRepo := &mockCourseRepo{}
	userRepo := &mockUserRepo{teachers: map[string]*model.User{}}
	service := NewCourseService(courseRepo, userRepo)
	input := CreateCourseInput{Name: "A", Remark: "", StudentMaxNum: 10, Hours: 2, StartDate: "2025-11-01", Semester: ""}
	_, _, err := service.CreateCourse("notfound", input)
	if err == nil {
		t.Error("should fail if teacher not found")
	}
}

func TestCreateCourse_SemesterValidate(t *testing.T) {
	courseRepo := &mockCourseRepo{}
	userRepo := &mockUserRepo{teachers: map[string]*model.User{"t1": {IDCard: "t1", Role: "teacher"}}}
	service := NewCourseService(courseRepo, userRepo)
	input := CreateCourseInput{Name: "A", Remark: "", StudentMaxNum: 10, Hours: 2, StartDate: "2026-03-07", Semester: "2026—11"}
	_, _, err := service.CreateCourse("t1", input)
	if err == nil {
		t.Error("should fail for invalid semester")
	}
}

func TestUpdateCourse_StudentNumLimit(t *testing.T) {
	courseRepo := &mockCourseRepo{courses: map[int64]*model.Course{1: {ID: 1, TeacherID: "t1"}}, getEnrollmentCount: 5}
	userRepo := &mockUserRepo{teachers: map[string]*model.User{"t1": {IDCard: "t1", Role: "teacher"}}}
	service := NewCourseService(courseRepo, userRepo)
	input := UpdateCourseInput{StudentMaxNum: intPtr(3)}
	_, _, err := service.UpdateCourse("t1", 1, input)
	if err == nil {
		t.Error("should fail if new student num < enrolled")
	}
}

func intPtr(i int) *int { return &i }
