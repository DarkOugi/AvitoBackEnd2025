-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS MerchUsers (
    id SERIAL PRIMARY KEY,
    merch_id VARCHAR(50) REFERENCES merch(name) NOT NULL ,
    user_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER NOT NULL CHECK ( cost > 0 ),
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS MerchUsers;
-- +goose StatementEnd
