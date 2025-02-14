package main

import (
	"avito/internal/db"
	"avito/internal/server"
	"avito/internal/service"
	"context"
	"github.com/fasthttp/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	var pSQL *db.PostgresDB

	ctx, stopSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.Stamp,
	}).Level(zerolog.DebugLevel)

	if pSQL == nil {
		pSQL, err = db.NewPostgresDB(ctx, "localhost", "5432", "avito", "0000", "avitodb")
		if err != nil {
			log.Fatal().Err(err).Msg("don't create connect with db")
		}
	}
	defer func() {
		err = pSQL.Close(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("db connection close error")
		}
	}()

	sv := service.NewService(pSQL)
	sr := server.NewServer(sv)

	r := router.New()
	r.POST("/api/auth", sr.Auth)
	r.GET("/api/buy/{item}", sr.BuyItem)
	r.POST("/api/sendCoin", sr.SendCoin)
	r.GET("/api/info", sr.Info)

	go func() {
		if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
			log.Fatal().Err(err).Msg("server critical error")
		}
	}()
	<-ctx.Done()
}
