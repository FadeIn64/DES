package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os/signal"
	"syscall"
)

var dbUrl = "postgres://username:password@localhost:5432/des"

func main() {
	seedCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	//init db
	dbConf, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(seedCtx, dbConf)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	//seeds
	err = seedLaps(seedCtx, pool)
	if err != nil {
		log.Fatal(err)
	}
}
