-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Merch (
    name VARCHAR(50) PRIMARY KEY,
    cost INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Merch;
-- +goose StatementEnd
