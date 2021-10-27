package repo

import (
	"github.com/ozonmp/lic-license-api/internal/model"
)

// TODO: Maybe LicenseRepo ?
type EventRepo interface {
	Lock(n uint64) ([]model.SubdomainEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.SubdomainEvent) error
	Remove(eventIDs []uint64) error
}
