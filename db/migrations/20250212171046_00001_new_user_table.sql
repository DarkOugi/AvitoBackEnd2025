-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users (
    login VARCHAR(50) PRIMARY KEY,
    password VARCHAR(512),
    balance INTEGER NOT NULL CHECK(balance >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd
