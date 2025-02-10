package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var dbConn *sql.DB = nil

// Здесь нужен царский синглтон чтоб м раз не конектиться к бд
func Connect(host, port, user, password, dbname string) (*sql.DB, error) {
	if dbConn != nil {
		return dbConn, nil
	}
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return dbConn, err
	} else {
		dbConn = db
	}
	// defer db.Close()

	return dbConn, err

}

// проверка бд
