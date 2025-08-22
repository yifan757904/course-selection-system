package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *model.Course) error
	GetByID(id int64) (*model.Course, error)
	GetByTeacherID(teacherID string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error)
	GetByTeacherName(teacherName string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error)
	GetByCourseName(courseName string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error)
	GetAll(pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error)
	Update(course *model.Course, updateData map[string]interface{}) error
	Delete(id int64) error
	GetEnrollmentCount(courseID int64) (int64, error)
}

type GormCourseRepository struct {
	db *gorm.DB
}

func NewGormCourseRepository(db *gorm.DB) *GormCourseRepository {
	return &GormCourseRepository{db: db}
}

func (r *GormCourseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *GormCourseRepository) GetByID(id int64) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	return &course, err
}

func (r *GormCourseRepository) GetByTeacherID(teacherID string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	var courses []map[string]interface{}
	var total int64

	query := r.db.Model(model.Course{})

	if err := query.Where("teacher_id = ? and semester = ?", teacherID, semester).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Order(sortBy + " " + sortOrder).
		Find(&courses).Error

	return courses, total, err
}

func (r *GormCourseRepository) GetByTeacherName(teacherName string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	var courses []map[string]interface{}
	var total int64

	// 先查询符合条件的教师
	var teachers []model.User
	if err := r.db.Where("name LIKE ? AND role = ? and semester = ?", "%"+teacherName+"%", "teacher", semester).Find(&teachers).Error; err != nil {
		return nil, 0, err
	}

	var teacherIDs []string
	for _, t := range teachers {
		teacherIDs = append(teacherIDs, t.IDCard)
	}

	if len(teacherIDs) == 0 {
		return nil, 0, nil
	}

	query := r.db.Model(model.Course{})
	if err := query.Where("teacher_id IN ?", teacherIDs).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Order(sortBy + " " + sortOrder).
		Find(&courses).Error

	return courses, total, err
}

func (r *GormCourseRepository) GetByCourseName(courseName string, pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	var courses []map[string]interface{}
	var total int64

	query := r.db.Model(&model.Course{})
	if err := query.Where("name LIKE ? and semester = ?", "%"+courseName+"%", semester).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Order(sortBy + " " + sortOrder).
		Find(&courses).Error

	return courses, total, err
}

func (r *GormCourseRepository) GetAll(pagination model.Pagination, sortBy, sortOrder, semester string) ([]map[string]interface{}, int64, error) {
	var courses []map[string]interface{}
	var total int64

	query := r.db.Model(&model.Course{}).Where("semester = ?", semester)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Order(sortBy + " " + sortOrder).
		Find(&courses).Error

	return courses, total, err
}

func (r *GormCourseRepository) Update(course *model.Course, updateData map[string]interface{}) error {
	return r.db.Model(course).Updates(updateData).Error
}

func (r *GormCourseRepository) Delete(id int64) error {
	return r.db.Delete(&model.Course{}, id).Error
}

func (r *GormCourseRepository) GetEnrollmentCount(courseID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Enrollment{}).Where("course_id = ?", courseID).Count(&count).Error
	return count, err
}
