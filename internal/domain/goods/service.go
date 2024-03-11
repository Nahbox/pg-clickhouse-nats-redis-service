package goods

import "context"

type Service interface {
	GetList(ctx context.Context, limit int, offset int) (*GetResponse, error)
	Create(data *Good) (*Good, error)
	Update(ctx context.Context, data *Good) (*Good, error)
	Remove(ctx context.Context, id int, projectId int) (*DeleteResponse, error)
	Reprioritize(ctx context.Context, id int, projectId, newPriority int) (*ReprioritizeResponse, error)
}
