-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS UserToUser (
    id SERIAL PRIMARY KEY,
    from_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    to_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER,
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX from_id_idx on UserToUser(from_id);
CREATE UNIQUE INDEX to_id_idx on UserToUser(to_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS UserToUser;
-- +goose StatementEnd
