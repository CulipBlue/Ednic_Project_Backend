package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/CulipBlue/backend_ednic/internal/config"
)

func Open(cfg config.Config) (*sql.DB, error) {
	if cfg.DatabaseDSN == "" {
		return nil, errors.New("DATABASE_DSN is required")
	}

	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DatabaseConnMaxLifetimeMi) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
