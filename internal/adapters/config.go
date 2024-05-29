package adapters

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

type Config struct {
	PostgresConnectionURI string `yaml:"postgres-connection-uri"`
}

func NewConfig() (*Config, error) {
	var conf Config
	err := cleanenv.ReadConfig("dev-config.yaml", &conf)
	if err != nil {
		slog.Debug("reading in config", err)
		return nil, err
	}

	return &conf, nil
}
