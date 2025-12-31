-- +goose Up
-- +goose StatementBegin
-- Add unique constraint to external_id to prevent duplicate films from external sources
CREATE UNIQUE INDEX idx_films_external_id ON films (external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_films_external_id;
-- +goose StatementEnd
