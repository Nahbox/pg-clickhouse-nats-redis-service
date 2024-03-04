package goods

import (
	"context"
)

type Repository interface {
	Add(data *Good) (*Good, error)
	Get(limit, offset int) (GetResponse, error)
	Update(data *Good) (*Good, error)
	Delete(ctx context.Context) error
}
