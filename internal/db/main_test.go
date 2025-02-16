package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testDB *PostgresDB

func TestMain(m *testing.M) {
	var err error

	//out, err := exec.Command("/bin/sh", "../../db/initDB.sh").Output()
	//fmt.Printf(string(out))
	//if err != nil {
	//	fmt.Printf(string(out))
	//	fmt.Printf("%s", err.Error())
	//	return
	//}

	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16-alpine",
		postgres.WithDatabase("avitodb"),
		postgres.WithInitScripts(filepath.Join("..", "..", "db", "t", "init_db_for_test.sql")),
		postgres.WithUsername("avito"),
		postgres.WithPassword("0000"),
		testcontainers.WithHostPortAccess(5432),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}
	connStr, err := postgresContainer.ConnectionString(context.Background(), "sslmode=disable", "application_name=test")

	conn, err := pgxpool.New(context.Background(), connStr)

	testDB = &PostgresDB{
		conn: conn,
	}
	//pgx.Connect(ctx, strСonn)
	if err != nil {

		return
	}
	//db, err := NewPostgresDB(context.Background(), "localhost", "5432", "avito", "0000", "avitodb")

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,
	}).Level(zerolog.DebugLevel)
	if err != nil {
		if err != nil {
			log.Fatal().Err(err).Msg("don't create connect with db")
		}
	}

	//testDB = db
	defer func() {
		testDB.Close()
	}()

	os.Exit(m.Run())
}
