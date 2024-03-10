package logs

import (
	"context"
)

type Repository interface {
	AddBatch(ctx context.Context, logs []Log) error
}
