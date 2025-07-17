package database

import (
	"database/sql"
	"fmt"
	"gin-demo/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func InitMySQL(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	return db, nil
}
