package repository

import (
	"context"
	"database/sql"
	"gin-demo/internal/model"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) CreateTeacher(ctx context.Context, teacher *model.Teacher) error {
	query := `
		INSERT INTO teacher (id, name)
		VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		teacher.ID,
		teacher.Name,
	)
	return err
}

func (r *TeacherRepository) GetTeacherByID(ctx context.Context, id string) (*model.Teacher, error) {
	query := `
		SELECT id, name
		FROM teacher WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var teacher model.Teacher
	err := row.Scan(
		&teacher.ID,
		&teacher.Name,
	)
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}
