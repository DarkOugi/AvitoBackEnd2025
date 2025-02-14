package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
}

func NewPostgresDB(ctx context.Context, host, port, username, password, dbname string) (*PostgresDB, error) {
	strСonn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, strСonn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) Close(ctx context.Context) error {
	err := db.conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("can't close connect: %w", err)
	}
	return nil
}
