package clickhouse

import (
	"context"
	_ "database/sql/driver"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chmigrate"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
)

var Migrations = chmigrate.NewMigrations()

func New(conf *config.CHConfig) (*ch.DB, error) {
	ctx := context.Background()
	dsn := conf.Dsn()

	db := ch.Connect(
		ch.WithDSN(dsn),
		ch.WithUser(conf.ChUser),
		ch.WithPassword(conf.ChPassword),
	)

	log.Info("clickhouse db connection established")

	err := Migrate(ctx, db, conf.CHMigrationsPath)
	if err != nil {
		return nil, err
	}

	log.Info("clickhouse db migrated")

	return db, nil
}

func Migrate(ctx context.Context, db *ch.DB, migrationsPath string) error {
	fsys := os.DirFS(migrationsPath)

	if err := Migrations.Discover(fsys); err != nil {
		return err
	}

	migrator := chmigrate.NewMigrator(db, Migrations)

	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	_, err = migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	return nil
}
