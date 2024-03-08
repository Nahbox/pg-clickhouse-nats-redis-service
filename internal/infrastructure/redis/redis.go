package redis

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
)

func New(conf *config.RConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: conf.RAddr(),
	})

	// Проверяем соединение с Redis
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Info("redis db connection established")

	return client, nil
}
