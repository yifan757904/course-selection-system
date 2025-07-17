package repository

import (
	"context"
	"database/sql"

	"github.com/liuyifan1996/internal/model"
)

type Course_StudentRepository struct {
	db *sql.DB
}

func NewCourse_StudentRepository(db *sql.DB) *Course_StudentRepository {
	return &Course_StudentRepository{db: db}
}

func (r *Course_StudentRepository) CreateCourse_Student(ctx context.Context, course_student *model.Course_Student) error {
	query := `
		INSERT INTO course_student (course_id, student_id)
		VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		course_student.CourseID,
		course_student.StudentID,
	)
	return err
}

func (r *Course_StudentRepository) GetNumByCourseID(ctx context.Context, id string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM course WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var num int
	err := row.Scan(
		num,
	)
	if err != nil {
		return -1, err
	}

	return num, nil
}
