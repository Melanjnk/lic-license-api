package repo

import (
	"github.com/jmoiron/sqlx"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
)

type LicenseEventRepo interface {
	Lock(n uint64) ([]model.LicenseEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.LicenseEvent) error // TODO: should trigger Created License Event?
	Remove(eventIDs []uint64) error
}

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) *eventRepo {
	return &eventRepo{db: db}
}

type eventPayload pb.License
