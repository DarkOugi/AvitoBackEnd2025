-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS MerchUsers (
    id SERIAL PRIMARY KEY,
    merch_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    cost INTEGER,
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS MerchUsers;
-- +goose StatementEnd
