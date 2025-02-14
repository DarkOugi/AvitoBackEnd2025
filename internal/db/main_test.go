package db

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
	"time"
)

var testDB *PostgresDB

func TestMain(m *testing.M) {
	var err error
	db, err := NewPostgresDB(context.Background(), "localhost", "5432", "avito", "0000", "avitodb")

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,
	}).Level(zerolog.DebugLevel)
	if err != nil {
		if err != nil {
			log.Fatal().Err(err).Msg("don't create connect with db")
		}
	}

	testDB = db
	defer func() {
		err = testDB.Close(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("db connection close error")
		}
	}()

	os.Exit(m.Run())
}
