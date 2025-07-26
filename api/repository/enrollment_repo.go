package repository

import (
	"github.com/liuyifan1996/course-selection-system/api/model"
	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) Create(enrollment *model.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *EnrollmentRepository) Delete(studentID, courseID int64) error {
	return r.db.Delete(&model.Enrollment{
		StudentID: studentID,
		CourseID:  courseID,
	}).Error
}

func (r *EnrollmentRepository) CountByCourse(courseID int64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Enrollment{}).
		Where("course_id = ?", courseID).
		Count(&count).Error
	return count, err
}
