package adapters

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PostgresConnectionURI string `yaml:"postgres-connection-uri" env:"POSTGRES_CONNECTION_URI" env-required:"true"`
	KafkaHost             string `yaml:"kafka-host" env:"KAFKA_HOST" env-required:"true"`
}

func NewConfig() (*Config, error) {
	var conf Config
	err := cleanenv.ReadConfig("dev-config.yaml", &conf)
	if err != nil {
		err = cleanenv.ReadEnv(&conf)
		if err != nil {
			return nil, err
		}
	}

	return &conf, nil
}
