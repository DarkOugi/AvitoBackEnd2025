package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	conn *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, host, port, username, password, dbname string) (*PostgresDB, error) {
	//nolint:nosprintfhostport // тише тише тише
	strConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)

	conn, err := pgxpool.New(ctx, strConn)
	// pgx.Connect(ctx, strСonn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) Close() {
	db.conn.Close()
}
