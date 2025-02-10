package initdb

// CREATE TABLES DB
// hash - index string

var Merch = `
	CREATE TABLE IF NOT EXISTS merch (
		name VARCHAR(50) PRIMARY KEY,
		cost INTEGER
	);
`
var Users = `
	CREATE TABLE IF NOT EXISTS users (
		login VARCHAR(50) PRIMARY KEY,
		password VARCHAR(50),
		balance INTEGER CHECK(balance >= 0)
	);
`

// добавить дату при создании строки
var MerchUser = `
	CREATE TABLE IF NOT EXISTS MerchUsers (
		id SERIAL PRIMARY KEY,
		merch_id VARCHAR(50) REFERENCES merch(name),
		user_id VARCHAR(50) REFERENCES users(login)
	);
`
var UserToUser = `
	CREATE TABLE IF NOT EXISTS UserToUser (
		id SERIAL PRIMARY KEY,
		from_id VARCHAR(50) REFERENCES users(login),
		to_id VARCHAR(50) REFERENCES users(login)
	);
`

// INSERT EXAMPLES
var InsertMerchValues = `
	INSERT INTO merch (name,cost)
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

// func Init_mech() (string, string) {
// 	return merch_init_table, merch_init_values
// }
