-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS MerchUsers (
    id SERIAL PRIMARY KEY,
    merch_id VARCHAR(50) REFERENCES merch(name) NOT NULL ,
    user_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER CHECK ( cost > 0 ),
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX march_id_idx on MerchUsers(merch_id);
CREATE UNIQUE INDEX user_id_idx on MerchUsers(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS MerchUsers;
-- +goose StatementEnd
