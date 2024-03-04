package domain

import "context"

type Repository interface {
	Add(ctx context.Context) error
	Get(ctx context.Context, limit, offset int) error
	Update(ctx context.Context) error
	Delete(ctx context.Context) error
}
