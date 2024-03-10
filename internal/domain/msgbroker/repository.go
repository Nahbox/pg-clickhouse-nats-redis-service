package msgbroker

import (
	"context"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/logs"
)

type Repository interface {
	Publish(data *logs.Log) error
	ReadAsync(ctx context.Context, subjectID string, out chan<- *logs.Log) error
}
