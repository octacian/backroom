-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS record (
	uuid char(27) NOT NULL PRIMARY KEY,
	cage VARCHAR(255) NOT NULL,
	data JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS record_cage ON record (cage);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS record;

DROP INDEX IF EXISTS record_cage;

-- +goose StatementEnd
