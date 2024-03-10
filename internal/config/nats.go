package config

import (
	"fmt"
)

type NatsConfig struct {
	Host string `envconfig:"NATS_HOST" required:"true"`
	Port int    `envconfig:"NATS_PORT" required:"true"`
}

func (r *NatsConfig) URL() string {
	return fmt.Sprintf("nats://%s:%d",
		r.Host, r.Port)
}
