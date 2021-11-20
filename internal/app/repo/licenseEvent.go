package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	model "github.com/ozonmp/lic-license-api/internal/model/license"
	pb "github.com/ozonmp/lic-license-api/pkg/lic-license-api"
	"time"
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

type licenseEvent struct {
	ID        uint64            `db:"id"`
	LicenseID uint64            `db:"license_id"`
	Type      model.EventType   `db:"type"`
	Status    model.EventStatus `db:"status"`
	License   eventPayload      `db:"payload"`
	UpdatedAt time.Time         `db:"updated"`
}

func (e *licenseEvent) toModel() model.LicenseEvent {
	return model.LicenseEvent{
		ID:        e.ID,
		LicenseID: e.LicenseID,
		Type:      e.Type,
		Status:    e.Status,
		License: &model.License{
			ID:    e.License.LicenseId,
			Title: e.License.Title,
		},
		UpdatedAt: e.UpdatedAt,
	}
}

func NewEventRepo(db *sqlx.DB) *eventRepo {
	return &eventRepo{db: db}
}

type eventPayload pb.License

func (r *eventRepo) Lock(ctx context.Context, n uint64) ([]model.LicenseEvent, error) {
	sqlStr := `
WITH lock_sellers AS (
    SELECT DISTINCT license_id FROM license_events WHERE status = 'LOCK'
)
UPDATE license_events e SET status = 'LOCK', updated = $1
FROM (
    SELECT e.id, e.license_id, e.type, e.status, e.payload, e.updated
    FROM license_events e
        LEFT JOIN license_events l ON (e.license_id = l.license_id)
    WHERE e.status = 'UNLOCK' AND l.license_id IS NULL
    ORDER BY e.id
    LIMIT $2
) AS q
WHERE e.id = q.id
RETURNING q.id, q.license_id, q.type, q.status, q.payload, q.updated
`
	rows, err := r.db.QueryContext(ctx, sqlStr, time.Now().UTC(), n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*licenseEvent, 0)
	err = sqlx.StructScan(rows, &events)
	if err != nil {
		return nil, err
	}

	res := make([]model.LicenseEvent, 0, len(events))
	for _, ev := range events {
		res = append(res, ev.toModel())
	}

	return res, nil
}

func (r *eventRepo) Unlock(ctx context.Context, eventIDs []uint64) error {
	q := sq.Update("license_events").PlaceholderFormat(sq.Dollar).
		Set("status", model.EventUnlock).
		Set("updated", time.Now().UTC()).
		Where(sq.Eq{"id": eventIDs})

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) Remove(ctx context.Context, eventIDs []uint64) error {
	q := sq.Delete("license_events").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": eventIDs})

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}
