package postgres

import (
	"context"
	"database/sql"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain"
)

type Repo struct {
	storage *sql.DB
}

func New(storage *sql.DB) domain.Repository {
	return &Repo{storage}
}

func (r *Repo) Add(ctx context.Context) error {
	return nil
}

func (r *Repo) Get(ctx context.Context) error {
	return nil
}

func (r *Repo) Update(ctx context.Context) error {
	return nil
}

func (r *Repo) Delete(ctx context.Context) error {
	return nil
}
