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
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Info("Подключение к Redis установлено:", pong)

	return client, nil
}
