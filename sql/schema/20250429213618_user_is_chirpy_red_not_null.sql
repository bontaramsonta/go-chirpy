-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN is_chirpy_red
SET
    NOT NULL;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN is_chirpy_red
DROP NOT NULL;

-- +goose StatementEnd
