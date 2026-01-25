-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN password_reset_token_hash TEXT,
ADD COLUMN password_reset_expire_at TIMESTAMP;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
DROP COLUMN password_reset_token_hash,
DROP COLUMN password_reset_expire_at;
-- +goose StatementEnd
