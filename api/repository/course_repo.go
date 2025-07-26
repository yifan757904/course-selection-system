package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *CourseRepository) FindByID(id int64) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	return &course, err
}

func (r *CourseRepository) FindAll() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) DeleteByTeacher(courseID, teacherID int64) error {
	result := r.db.Where("id = ? AND teacher_id = ?", courseID, teacherID).Delete(&model.Course{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CourseRepository) FindByTeacher(teacherID int64) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Where("teacher_id = ?", teacherID).Find(&courses).Error
	return courses, err
}
