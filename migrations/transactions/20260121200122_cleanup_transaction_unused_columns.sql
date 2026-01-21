-- +goose Up
-- +goose StatementBegin
ALTER TABLE transaction
  DROP COLUMN IF EXISTS refunded_at;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE transaction
  ADD COLUMN IF NOT EXISTS refunded_at timestamptz;
-- +goose StatementEnd
