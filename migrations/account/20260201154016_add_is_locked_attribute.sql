-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_account
ADD COLUMN is_locked BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN lockout_counter INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_account
DROP COLUMN IF EXISTS is_locked,
DROP COLUMN IF EXISTS lockout_counter;
-- +goose StatementEnd
