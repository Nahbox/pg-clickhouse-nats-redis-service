package clickhouse

import (
	"database/sql/driver"
	"errors"
	"path/filepath"

	"github.com/ClickHouse/clickhouse-go"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New() (driver.Conn, error) {
	connStr := "tcp://localhost:9000?username=default&password="

	// Подключение к ClickHouse
	chConn, err := clickhouse.Open(connStr)
	if err != nil {
		return nil, err
	}

	err = Migrate(connStr)
	if err != nil {
		return nil, err
	}

	return chConn, nil
}

func Migrate(connStr string) error {
	sourceUrl := "file://" + filepath.Join("/", "internal", "db", "clickhouse", "migrations")

	m, err := migrate.New(sourceUrl, connStr)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
