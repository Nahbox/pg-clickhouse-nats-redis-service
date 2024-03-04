package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	model "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
)

type Repo struct {
	storage *sql.DB
}

func New(storage *sql.DB) model.Repository {
	return &Repo{storage}
}

func (r *Repo) Add(data *model.Good) (*model.Good, error) {
	data.Removed = false

	tx, err := r.storage.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
			INSERT INTO goods (project_id, name,  description, removed) 
			VALUES ($1, $2, $3, $4, $5) RETURNING id, priority, created_at;
			`,
		data.ProjectId, data.Name, data.Description, data.Removed).Scan(&data.Id, &data.Priority, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return data, tx.Commit()
}

// TODO: dodelat'
func (r *Repo) Get(limit, offset int) (model.GetResponse, error) {
	if limit == 0 {
		limit = 10 // Устанавливаем значение по умолчанию
	}

	var goods []model.Good

	rows, err := r.storage.Query("SELECT *, COUNT(*) AS total, FROM goods LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return model.GetResponse{Meta: model.Meta{Total: 0, Removed: 0, Limit: limit, Offset: offset}, Goods: goods}, err
	}

	total := 0
	removedCount := 0

	for rows.Next() {
		var good model.Good
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt, &total)
		if err != nil {
			return model.GetResponse{Meta: model.Meta{Total: 0, Removed: 0, Limit: limit, Offset: offset}, Goods: goods}, err
		}
		goods = append(goods, good)
	}

	for _, good := range goods {
		if good.Removed == true {
			removedCount++
		}
	}

	metaInfo := model.Meta{Total: total, Removed: removedCount, Limit: limit, Offset: offset}

	res := model.GetResponse{Meta: metaInfo, Goods: goods}

	return res, nil
}

func (r *Repo) Update(data *model.Good) (*model.Good, error) {
	data.Removed = false

	tx, err := r.storage.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
			UPDATE goods SET name=$1, description=$2 WHERE id=$3, project_id=$4 RETURNING priority, removed, created_at;
			`,
		data.Name, data.Description, data.Id, data.ProjectId).Scan(&data.Priority, &data.Removed, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return data, tx.Commit()
}

func (r *Repo) Delete(ctx context.Context) error {

	return nil
}
