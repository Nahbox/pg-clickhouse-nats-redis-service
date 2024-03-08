package clickhouse

import (
	"database/sql/driver"
	"errors"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
)

func New(conf *config.CHConfig) (driver.Conn, error) {
	connStr := conf.ChDsn()
	sourceUrl := conf.ChMigrationsPathStr()

	// Подключение к ClickHouse
	chConn, err := clickhouse.Open(connStr)
	if err != nil {
		return nil, err
	}

	err = Migrate(connStr, sourceUrl)
	if err != nil {
		return nil, err
	}

	return chConn, nil
}

func Migrate(connStr string, sourceUrl string) error {
	m, err := migrate.New(sourceUrl, connStr)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
