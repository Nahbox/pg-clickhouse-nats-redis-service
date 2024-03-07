package kvstore

import (
	"context"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
)

type Repository interface {
	Set(ctx context.Context, key string, value interface{}) error
	GetList(ctx context.Context, key string) (*goods.GetResponse, error)
}
