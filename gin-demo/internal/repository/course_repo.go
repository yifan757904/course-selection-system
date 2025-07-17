package repository

import (
	"context"
	"database/sql"
	"gin-demo/internal/model"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) CreateCourse(ctx context.Context, course *model.Course) error {
	query := `
		INSERT INTO course (id, name, teacher_id, teacher_name, remarks, student_maxnum, time_max, time_min)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.Name,
		course.TeacherID,
		course.TeacherName,
		course.Remarks,
		course.Student_maxnum,
		course.Time_max,
		course.Time_min,
	)
	return err
}

func (r *CourseRepository) GetCourseByID(ctx context.Context, id string) (*model.Course, error) {
	query := `
		SELECT *
		FROM course WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var course model.Course
	err := row.Scan(
		&course.ID,
		&course.Name,
		&course.TeacherID,
		&course.TeacherName,
		&course.Remarks,
		&course.Student_maxnum,
		&course.Time_max,
		&course.Time_min,
	)
	if err != nil {
		return nil, err
	}

	return &course, nil
}
