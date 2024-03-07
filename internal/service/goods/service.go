package goods

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/kvstore"
)

type Service struct {
	goodsRepo   goods.Repository
	kvstoreRepo kvstore.Repository
}

func NewService(goodsRepo goods.Repository, kvstoreRepo kvstore.Repository) *Service {
	return &Service{goodsRepo, kvstoreRepo}
}

func (s *Service) GetList(ctx context.Context, limit int, offset int) (*goods.GetResponse, error) {
	var res *goods.GetResponse

	// Запрос в redis
	key := fmt.Sprintf("l=%s,o=%s", strconv.Itoa(limit), strconv.Itoa(offset))
	res, err := s.kvstoreRepo.GetList(ctx, key)

	// Если запись найдена, возвращаем ее
	if res != nil && err == nil {
		return res, err
	}

	// Иначе возвращаем результат из postgres
	res, err = s.goodsRepo.Get(limit, offset)
	if res != nil && err != nil {
		key := fmt.Sprintf("o=%s,l=%s", strconv.Itoa(res.Meta.Offset), strconv.Itoa(res.Meta.Limit))

		// Помещаем результат в redis
		err = s.kvstoreRepo.Set(ctx, key, res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (s *Service) Create(data *goods.Good) (*goods.Good, error) {
	return s.goodsRepo.Add(data)
}

func (s *Service) Update(data *goods.Good) (*goods.Good, error) {
	return s.goodsRepo.Update(data)
}

func (s *Service) Remove(id int, projectId int) (*goods.DeleteResponse, error) {
	return s.goodsRepo.Delete(id, projectId)
}

func (s *Service) Reprioritize(id int, projectId, newPriority int) (*goods.ReprioritizeResponse, error) {
	return s.goodsRepo.UpdatePriority(id, projectId, newPriority)
}
