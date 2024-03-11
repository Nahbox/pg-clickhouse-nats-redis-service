package goods

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/kvstore"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/logs"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/msgbroker"
)

type Service struct {
	goodsRepo   goods.Repository
	kvstoreRepo kvstore.Repository
	msgbRepo    msgbroker.Repository
}

func NewService(goodsRepo goods.Repository, kvstoreRepo kvstore.Repository, msgbRepo msgbroker.Repository) *Service {
	return &Service{goodsRepo, kvstoreRepo, msgbRepo}
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
	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("l=%s,o=%s", strconv.Itoa(res.Meta.Limit), strconv.Itoa(res.Meta.Offset))
	// Помещаем результат в redis
	err = s.kvstoreRepo.Set(ctx, key, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *Service) Create(data *goods.Good) (*goods.Good, error) {
	resp, err := s.goodsRepo.Add(data)
	if err != nil {
		return nil, err
	}

	logData := goodToLog(data)
	err = s.msgbRepo.Publish(logData)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, data *goods.Good) (*goods.Good, error) {
	resp, err := s.goodsRepo.Update(data)
	if err != nil {
		return nil, err
	}

	err = s.kvstoreRepo.EraseAll(ctx)
	if err != nil {
		return nil, err
	}

	logData := goodToLog(data)
	err = s.msgbRepo.Publish(logData)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Remove(ctx context.Context, id int, projectId int) (*goods.DeleteResponse, error) {
	resp, good, err := s.goodsRepo.Delete(id, projectId)
	if err != nil {
		return nil, err
	}

	err = s.kvstoreRepo.EraseAll(ctx)
	if err != nil {
		return nil, err
	}

	logData := goodToLog(good)
	err = s.msgbRepo.Publish(logData)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Reprioritize(ctx context.Context, id int, projectId, newPriority int) (*goods.ReprioritizeResponse, error) {
	resp, goodsData, err := s.goodsRepo.UpdatePriority(id, projectId, newPriority)
	if err != nil {
		return nil, err
	}

	err = s.kvstoreRepo.EraseAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, good := range goodsData {
		logData := goodToLog(&good)
		err = s.msgbRepo.Publish(logData)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func goodToLog(goodData *goods.Good) *logs.Log {
	return &logs.Log{
		Id:          goodData.Id,
		ProjectId:   goodData.ProjectId,
		Name:        goodData.Name,
		Description: goodData.Description,
		Priority:    goodData.Priority,
		Removed:     goodData.Removed,
		EventTime:   goodData.CreatedAt,
	}
}
