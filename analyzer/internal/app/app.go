package app

import (
	"DAS/internal/consumers"
	"DAS/internal/controllers"
	"DAS/internal/metrics"
	"DAS/internal/web"
	"context"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"

	"DAS/config"
	"DAS/internal/repositories"
)

type App struct {
	Cfg       *config.Config
	pool      *pgxpool.Pool
	trManager trm.Manager

	lapRepo     *repositories.LapRepository
	driverRepo  *repositories.DriverRepository
	teamRepo    *repositories.TeamRepository
	meetingRepo *repositories.MeetingRepository

	LapHandler *consumers.LapHandler
	Exporter   *metrics.Exporter

	Server web.HttpServer
}

func NewApp(cfg *config.Config) *App {
	pool := initDBPool(cfg)
	trManager := initTransactionManager(pool)

	lapRepo := repositories.NewLapRepository(pool, trManager)
	driverRepo := repositories.NewDriverRepository(pool, trManager)
	teamRepo := repositories.NewTeamRepository(pool, trManager)
	meetingRepo := repositories.NewMeetingRepository(pool, trManager)

	meetingCtrl := controllers.NewMeetingController(meetingRepo)
	driverCtrl := controllers.NewDriverController(driverRepo)
	teamCtrl := controllers.NewTeamController(teamRepo)

	exporter := metrics.NewMetricsExporter()
	lapHandler := consumers.NewLapHandler(lapRepo, exporter)

	db := stdlib.OpenDBFromPool(pool)

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatal(err)
	}

	server, err := web.New(cfg.ServerPort, &metricsHandler{pronH: promhttp.Handler()},
		meetingCtrl,
		driverCtrl,
		teamCtrl,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		Cfg:         cfg,
		pool:        pool,
		trManager:   trManager,
		lapRepo:     lapRepo,
		driverRepo:  driverRepo,
		teamRepo:    teamRepo,
		meetingRepo: meetingRepo,
		LapHandler:  lapHandler,
		Exporter:    exporter,
		Server:      server,
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
	)
}

func (a *App) Close() {
	a.pool.Close()
}

type metricsHandler struct {
	pronH http.Handler
}

func (h metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving metrics at %s", r.URL.Path)
	h.pronH.ServeHTTP(w, r)
}
