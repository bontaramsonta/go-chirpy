-- +goose Up
-- +goose StatementBegin
ALTER TABLE chirps
RENAME COLUMN content TO body;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE chirps
RENAME COLUMN body TO content;

-- +goose StatementEnd
