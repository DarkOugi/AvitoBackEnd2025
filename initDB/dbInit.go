package initdb

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// CREATE TABLES DB

var Merch = `
	CREATE TABLE IF NOT EXISTS Merch (
		name VARCHAR(50) PRIMARY KEY,
		cost INTEGER
	);
`
var MerchIndex = `
	CREATE INDEX IF NOT EXISTS idx_hash_index ON Merch USING HASH (name); 
`
var Users = `
	CREATE TABLE IF NOT EXISTS Users (
		login VARCHAR(50) PRIMARY KEY,
		password VARCHAR(512),
		balance INTEGER CHECK(balance >= 0)
	);
`
var UsersIndex = `
	CREATE INDEX IF NOT EXISTS idx_hash_index ON Users USING HASH (login);
`

// добавить дату при создании строки
var MerchUser = `
	CREATE TABLE IF NOT EXISTS MerchUsers (
		id SERIAL PRIMARY KEY,
		merch_id VARCHAR(50) NOT NULL,
	    user_id VARCHAR(50) NOT NULL,
	    cost INTEGER,
	    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
var MerchUIndex = `
	CREATE INDEX IF NOT EXISTS idx_hash_index ON MerchUsers USING HASH (merch_id);
`
var UserMIndex = `
	CREATE INDEX IF NOT EXISTS  idx_hash_index ON MerchUsers USING HASH (user_id);
`
var UserToUser = `
	CREATE TABLE IF NOT EXISTS UserToUser (
		id SERIAL PRIMARY KEY,
		from_id VARCHAR(50) NOT NULL,
	    to_id VARCHAR(50) NOT NULL,
	    cost INTEGER,
	    transactionTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

var FromIndex = `
	CREATE INDEX IF NOT EXISTS idx_hash_index ON UserToUser USING HASH (from_id);
`
var ToIndex = `
	CREATE INDEX IF NOT EXISTS idx_hash_index ON UserToUsers USING HASH (to_id);
`

// INSERT EXAMPLES
var insertMerchValues = `
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
`
var tables = []string{
	Users,
	//Merch,
	//UserToUser,
	//MerchUser,
}
var indexes = []string{UsersIndex, MerchIndex, MerchUIndex, UserMIndex, FromIndex, ToIndex}

func CreateTables(db *pgx.Conn) error {
	for _, table := range tables {
		_, err := db.Exec(context.Background(), table)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateIndex(db *pgx.Conn) error {
	for _, index := range indexes {
		_, err := db.Exec(context.Background(), index)
		if err != nil {
			return err
		}
	}
	return nil
}
func InsertValue(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), insertMerchValues)
	if err != nil {
		return err
	}
	return nil
}
