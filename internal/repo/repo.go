package repo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/lic-license-api/internal/model"
)

// Repo is DAO for Template
type Repo interface {
	CreateLicense(ctx context.Context, license *model.License) (uint64, error)
	DescribeLicense(ctx context.Context, licenseID uint64) (*model.License, error)
	ListLicense(ctx context.Context, cursor uint64, limit uint64) ([]*model.License, error)
	RemoveLicense(ctx context.Context, licenseID uint64) (bool, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) CreateLicense(ctx context.Context, license *model.License) (uint64, error) {
	return 0, errors.New("not implemented")
}

func (r *repo) DescribeLicense(ctx context.Context, licenseID uint64) (*model.License, error) {
	return nil, errors.New("not implemented")
}

func (r *repo) ListLicense(ctx context.Context, cursor uint64, limit uint64) ([]*model.License, error) {
	return nil, errors.New("not implemented")
}

func (r *repo) RemoveLicense(ctx context.Context, licenseID uint64) (bool, error) {
	return false, errors.New("not implemented")
}
