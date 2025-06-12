package config

import (
	"os"
	"strings"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/spf13/viper"
)

const EnvPrefix = "DAS"

type Config struct {
	KafkaBroker  string `mapstructure:"kafka_broker"`
	KafkaTopic   string `mapstructure:"kafka_topic"`
	KafkaGroupID string `mapstructure:"kafka_group_id"`
	PGConnString string `mapstructure:"db_url"`
	ServerPort   string `mapstructure:"server_port"`
}

func (a *Config) TransactionSettings() settings.Settings {
	return settings.Must(
		settings.WithTimeout(time.Second * 5),
	)
}

func Get() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath("./")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		v.Set(k, os.ExpandEnv(val))
	}

	var cfg *Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg, nil
}
