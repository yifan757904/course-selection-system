package repository

import (
	"context"
	"database/sql"
	"gin-demo/internal/model"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, stu *model.Student) error {
	query := `
		INSERT INTO student (id, name)
		VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		stu.ID,
		stu.Name,
	)
	return err
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	query := `
		SELECT id, name
		FROM student WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var student model.Student
	err := row.Scan(
		&student.ID,
		&student.Name,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}
