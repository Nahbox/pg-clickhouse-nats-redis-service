package config

import "fmt"

type CHConfig struct {
	ChHost           string `envconfig:"CLICKHOUSE_HOST" required:"true"`
	ChPort           int    `envconfig:"CLICKHOUSE_PORT" required:"true"`
	ChUser           string `envconfig:"CLICKHOUSE_USER" required:"true"`
	ChPassword       string `envconfig:"CLICKHOUSE_PASSWORD" required:"true"`
	CHMigrationsPath string `envconfig:"CLICKHOUSE_MIGRATIONS_PATH" required:"true"`
}

func (ch *CHConfig) ChDsn() string {
	return fmt.Sprintf("clickhouse://%s:%d?username=%s&password=%s",
		ch.ChHost, ch.ChPort, ch.ChUser, ch.ChPassword)
}

func (ch *CHConfig) ChMigrationsPathStr() string {
	return fmt.Sprintf("file://%s",
		ch.CHMigrationsPath)
}
