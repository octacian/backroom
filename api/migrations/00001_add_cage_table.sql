-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cage (
	uuid char(27) NOT NULL PRIMARY KEY,
	key VARCHAR(255) NOT NULL,
	data JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS cage_key ON cage (key);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cage;

DROP INDEX IF EXISTS cage_key;

-- +goose StatementEnd
