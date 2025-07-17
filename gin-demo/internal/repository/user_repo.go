package repository

import (
	"context"
	"database/sql"
	"gin-demo/internal/model"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, password, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Password,
		user.Email,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password, email, created_at, updated_at
		FROM users WHERE username = ?
	`
	row := r.db.QueryRowContext(ctx, query, username)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
