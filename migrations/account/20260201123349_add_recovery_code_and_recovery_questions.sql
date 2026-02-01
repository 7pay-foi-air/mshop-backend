-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_account
ADD COLUMN recovery_code_location text,
ADD COLUMN security_questions_hash text;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_account
DROP COLUMN IF EXISTS recovery_code_location,
DROP COLUMN IF EXISTS security_questions_hash;
-- +goose StatementEnd