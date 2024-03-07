package goods

import "context"

type Service interface {
	GetList(ctx context.Context, limit int, offset int) (*GetResponse, error)
	Create(data *Good) (*Good, error)
	Update(data *Good) (*Good, error)
	Remove(id int, projectId int) (*DeleteResponse, error)
	Reprioritize(id int, projectId, newPriority int) (*ReprioritizeResponse, error)
}
