package repo

import (
	"github.com/ozonmp/lic-license-api/internal/model"
)

// TODO: Think about is it Event?
type LicenseEventRepo interface {
	Lock(n uint64) ([]model.LicenseEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.LicenseEvent) error // TODO: should trigger Created License Event?
	Remove(eventIDs []uint64) error
}
