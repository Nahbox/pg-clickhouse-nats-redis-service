package clickhouse

import (
	"context"

	"github.com/uptrace/go-clickhouse/ch"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/logs"
)

type Repo struct {
	db *ch.DB
}

func NewRepo(db *ch.DB) logs.Repository {
	return &Repo{db}
}

func (r *Repo) AddBatch(ctx context.Context, batch []logs.Log) error {
	_, err := r.db.NewInsert().Model(&batch).Exec(ctx)
	return err
}
