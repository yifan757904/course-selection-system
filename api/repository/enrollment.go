package repository

import (
	"fmt"

	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

// 便于测试的接口抽象
type EnrollmentRepo interface {
	GetStudentByIDCard(idCard string) (*model.User, error)
	GetCourseByID(courseID int) (*model.Course, error)
	GetEnrollment(studentID, courseID int64) (*model.Enrollment, error)
	CreateEnrollment(enrollment *model.Enrollment) error
	DeleteEnrollment(enrollment *model.Enrollment) error
	CountEnrollmentsByCourse(courseID int64) (int64, error)
	GetStudentEnrollments(studentID int64) ([]model.Enrollment, error)
	GetStudentsByCourseID(courseID int64) ([]model.User, error)
	GetStudentCourses(enrollmentIDs []int64, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error)
}

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) GetStudentByIDCard(idCard string) (*model.User, error) {
	var student model.User
	err := r.db.Where("id_card = ? AND role = 'student'", idCard).First(&student).Error
	return &student, err
}

func (r *EnrollmentRepository) GetCourseByID(courseID int) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, courseID).Error
	return &course, err
}

func (r *EnrollmentRepository) GetEnrollment(studentID, courseID int64) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&enrollment).Error
	return &enrollment, err
}

func (r *EnrollmentRepository) CreateEnrollment(enrollment *model.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *EnrollmentRepository) DeleteEnrollment(enrollment *model.Enrollment) error {
	return r.db.Delete(enrollment).Error
}

func (r *EnrollmentRepository) CountEnrollmentsByCourse(courseID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}

func (r *EnrollmentRepository) GetStudentEnrollments(studentID int64) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	err := r.db.Where("student_id = ?", studentID).Find(&enrollments).Error
	return enrollments, err
}
func (r *EnrollmentRepository) GetStudentsByCourseID(courseID int64) ([]model.User, error) {
	var students []model.User
	err := r.db.Table("users").
		Select("users.*").
		Joins("JOIN enrollments ON enrollments.student_id = users.id").
		Where("enrollments.course_id = ? AND users.role = ?", courseID, "student").
		Find(&students).Error
	return students, err
}

func (r *EnrollmentRepository) GetStudentCourses(enrollmentIDs []int64, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	var courses []map[string]interface{}
	var total int64

	query := r.db.Model(&model.Course{})
	query = query.Where("id IN ?", enrollmentIDs)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&courses).Error

	return courses, total, err
}
