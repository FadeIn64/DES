package config

import (
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

type App struct {
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	PGConnString string
	SectorsCount int
}

func (a *App) TransactionSettings() settings.Settings {
	return settings.Must(
		settings.WithTimeout(time.Second * 5),
	)
}
