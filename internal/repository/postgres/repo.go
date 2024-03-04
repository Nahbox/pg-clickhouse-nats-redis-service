package postgres

import (
	"context"
	"database/sql"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain"
	pgModel "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/models/postgres"
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

func (r *Repo) Get(ctx context.Context, limit, offset int) error {
	var projects []pgModel.Projects
	var goods []pgModel.Goods

	return nil
}

func (r *Repo) Update(ctx context.Context) error {
	return nil
}

func (r *Repo) Delete(ctx context.Context) error {
	return nil
}
