package app

import (
	"context"
	"log"

	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	"DAS/config"
	"DAS/internal/repositories"
)

type App struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	trManager trm.Manager
	Repo      *repositories.LapRepository
}

func NewApp(cfg *config.Config) *App {
	pool := initDBPool(cfg)
	trManager := initTransactionManager(pool)
	repo := repositories.NewLapRepository(pool, trManager)

	return &App{
		cfg:       cfg,
		pool:      pool,
		trManager: trManager,
		Repo:      repo,
	}
}

func initDBPool(cfg *config.Config) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(cfg.PGConnString)
	if err != nil {
		log.Fatalf("Unable to parse pool config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}

func initTransactionManager(pool *pgxpool.Pool) trm.Manager {

	f := trmpgx.NewDefaultFactory(pool)

	return manager.Must(
		f,
		//manager.WithSettings(*config.NewApp().TransactionSettings()),
	)
}

func (a *App) Close() {
	a.pool.Close()
}
