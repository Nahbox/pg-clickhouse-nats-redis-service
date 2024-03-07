package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort  int `envconfig:"APP_PORT" required:"true"`
	PgConfig *PgConfig
	CHConfig *CHConfig
	RConfig  *RConfig
}

func FromEnv() (*Config, error) {
	cfg := Config{}

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
