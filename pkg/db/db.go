package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Connect(host, port, username, password, dbname string) (*pgx.Conn, error) {
	str_conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), str_conn)
	return conn, err
}
