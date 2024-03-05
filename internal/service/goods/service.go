package goods

import (
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
)

type Service struct {
	repo goods.Repository
}

func NewService(repo goods.Repository) (*Service, error) {
	return &Service{repo}, nil
}

func (s *Service) GetList(limit int, offset int) (*goods.GetResponse, error) {
	return s.repo.Get(limit, offset)
}

func (s *Service) Create(data *goods.Good) (*goods.Good, error) {
	return s.repo.Add(data)
}

func (s *Service) Update(data *goods.Good) (*goods.Good, error) {
	return s.repo.Update(data)
}

func (s *Service) Remove(id int, projectId int) (*goods.DeleteResponse, error) {
	return s.repo.Delete(id, projectId)
}

func (s *Service) Reprioritize(id int, projectId, newPriority int) (*goods.ReprioritizeResponse, error) {
	return s.repo.UpdatePriority(id, projectId, newPriority)
}
