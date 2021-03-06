package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/lic-license-api/internal/pkg/logger"
)

// NewPostgres returns DB
func NewPostgres(ctx context.Context, dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		logger.ErrorKV(ctx, "Failed to create database connection", "err", err)

		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.ErrorKV(ctx, "Failed ping the database", "err", err)

		return nil, err
	}

	return db, nil
}
