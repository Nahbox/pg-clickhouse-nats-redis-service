package config

import "fmt"

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	User     string `envconfig:"REDIS_USER" required:"true"`
	Password string `envconfig:"REDIS_PASSWORD" required:"true"`
	Port     int    `envconfig:"REDIS_PORT" required:"true"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d",
		r.Host, r.Port)
}
