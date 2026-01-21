-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_account
  DROP COLUMN IF EXISTS is_email_verified,
  DROP COLUMN IF EXISTS email_verification_token_hash,
  DROP COLUMN IF EXISTS email_verification_expires_at;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_account
  ADD COLUMN IF NOT EXISTS is_email_verified boolean NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS email_verification_token_hash text,
  ADD COLUMN IF NOT EXISTS email_verification_expires_at timestamptz;
-- +goose StatementEnd
