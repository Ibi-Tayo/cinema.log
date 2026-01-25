-- +goose Up
-- +goose StatementBegin
-- Make github_id nullable to support multiple auth providers
ALTER TABLE users ALTER COLUMN github_id DROP NOT NULL;

-- Add google_id column for Google OAuth support
ALTER TABLE users ADD COLUMN google_id VARCHAR(255);

-- Add unique constraint to google_id
ALTER TABLE users ADD CONSTRAINT users_google_id_unique UNIQUE (google_id);

-- Create index on google_id for faster lookups
CREATE INDEX idx_users_google_id ON users(google_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove index
DROP INDEX IF EXISTS idx_users_google_id;

-- Remove unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_google_id_unique;

-- Remove google_id column
ALTER TABLE users DROP COLUMN IF EXISTS google_id;

-- Restore github_id NOT NULL constraint
-- Note: This will fail if there are users without github_id
ALTER TABLE users ALTER COLUMN github_id SET NOT NULL;
-- +goose StatementEnd
