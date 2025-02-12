package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Connect(host, port, username, password, dbname string) (*pgx.Conn, error) {
	strСonn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), strСonn)
	return conn, err
}
