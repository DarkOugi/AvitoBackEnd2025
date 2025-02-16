CREATE TABLE IF NOT EXISTS Users (
     login VARCHAR(50) PRIMARY KEY,
     password VARCHAR(512),
     balance INTEGER NOT NULL CHECK(balance >= 0)
);
CREATE TABLE IF NOT EXISTS Merch (
     name VARCHAR(50) PRIMARY KEY,
     cost INTEGER NOT NULL CHECK(cost >= 0)
);

CREATE TABLE IF NOT EXISTS MerchUsers (
    id SERIAL PRIMARY KEY,
    merch_id VARCHAR(50) REFERENCES merch(name) NOT NULL ,
    user_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER NOT NULL CHECK ( cost > 0 ),
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS UserToUser (
    id SERIAL PRIMARY KEY,
    from_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    to_id VARCHAR(50) REFERENCES users(login) NOT NULL,
    cost INTEGER NOT NULL ,
    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO Merch (name,cost)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500)
;