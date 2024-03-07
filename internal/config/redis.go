package config

import "fmt"

type RConfig struct {
	RHost     string `envconfig:"REDIS_HOST" required:"true"`
	RUser     string `envconfig:"REDIS_USER" required:"true"`
	RPassword string `envconfig:"REDIS_PASSWORD" required:"true"`
	RPort     int    `envconfig:"REDIS_PORT" required:"true"`
}

func (r *RConfig) RAddr() string {
	return fmt.Sprintf("%s:%d",
		r.RHost, r.RPort)
}
