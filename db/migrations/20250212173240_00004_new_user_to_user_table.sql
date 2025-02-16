-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS UserToUser (
    id SERIAL PRIMARY KEY,
    from_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    to_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER NOT NULL ,
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS UserToUser;
-- +goose StatementEnd
