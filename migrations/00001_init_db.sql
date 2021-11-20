-- +goose Up
CREATE TABLE licenses
(
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR   NOT NULL,
    is_removed  BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE INDEX licenses_is_removed_idx ON licenses USING btree (is_removed);

CREATE TYPE license_event_type AS ENUM ('CREATED', 'UPDATED', 'REMOVED');
CREATE TYPE license_event_status AS ENUM ('DEFERRED', 'PROCESSED');

CREATE TABLE license_events
(
    id         BIGSERIAL PRIMARY KEY,
    license_id BIGSERIAL             NOT NULL,
    type       license_event_type    NOT NULL,
    status     license_event_status  NOT NULL,
    payload    JSONB                 NOT NULL,
    updated_at TIMESTAMP             NOT NULL
);

CREATE INDEX license_events_license_id_idx ON license_events USING btree (license_id);
CREATE INDEX license_events_type_idx ON license_events USING btree (type);
CREATE INDEX license_events_status_idx ON license_events USING btree (status);

-- +goose Down
DROP INDEX licenses_is_removed_idx;
DROP INDEX license_events_status_idx;
DROP INDEX license_events_type_idx;
DROP INDEX license_events_license_id_idx;
DROP TABLE licenses;
DROP TABLE license_events;
DROP TYPE license_event_type;
DROP TYPE license_event_status;
