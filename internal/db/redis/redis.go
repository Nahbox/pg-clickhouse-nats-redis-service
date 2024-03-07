package redis

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Проверяем соединение с Redis
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Info("Подключение к Redis установлено:", pong)

	return client, nil
}
