package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/kvstore"
)

type Repo struct {
	storage *redis.Client
}

func NewKVStoreRepo(storage *redis.Client) kvstore.Repository {
	return &Repo{storage}
}

func (r *Repo) Set(ctx context.Context, key string, value interface{}) error {
	// Сериализация значения в байтовый вид
	encodedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Запись данных в Redis с истекающим сроком годности (время жизни одна минута)
	err = r.storage.Set(ctx, key, encodedValue, time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetList(ctx context.Context, key string) (*goods.GetResponse, error) {
	var res *goods.GetResponse

	// Получаем значение для данного ключа
	val, err := r.storage.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repo) EraseAll(ctx context.Context) error {
	return r.storage.FlushDB(ctx).Err()
}
