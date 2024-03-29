package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	model "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/domain/goods"
)

type Repo struct {
	storage *sql.DB
}

func NewGoodsRepo(storage *sql.DB) model.Repository {
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
			VALUES ($1, $2, $3, $4) RETURNING id, priority, created_at;
			`,
		data.ProjectId, data.Name, data.Description, data.Removed).Scan(&data.Id, &data.Priority, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return data, tx.Commit()
}

func (r *Repo) Get(limit, offset int) (*model.GetResponse, error) {
	if limit == 0 {
		limit = 10 // Устанавливаем значение по умолчанию
	}

	var goods []model.Good

	rows, err := r.storage.Query(`
	SELECT *,
    (SELECT COUNT(*) FROM goods) AS total,
    (SELECT COUNT(*) FROM goods WHERE removed = false) AS removed_count
	FROM goods LIMIT $1 OFFSET $2;
	`,
		limit, offset)
	if err != nil {
		goods = []model.Good{} // Если записей нет, возвращаем пустой массив
		return &model.GetResponse{Meta: model.Meta{Total: 0, Removed: 0, Limit: limit, Offset: offset}, Goods: goods}, err
	}

	total := 0
	removedCount := 0

	for rows.Next() {
		var good model.Good
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt, &total, &removedCount)
		if err != nil {
			goods = []model.Good{} // Если записей нет, возвращаем пустой массив
			return &model.GetResponse{Meta: model.Meta{Total: 0, Removed: 0, Limit: limit, Offset: offset}, Goods: goods}, err
		}
		goods = append(goods, good)
	}
	defer rows.Close()

	if goods == nil {
		goods = []model.Good{} // Если записей нет, возвращаем пустой массив
	}

	metaInfo := model.Meta{Total: total, Removed: removedCount, Limit: limit, Offset: offset}

	res := &model.GetResponse{Meta: metaInfo, Goods: goods}

	return res, nil
}

func (r *Repo) Update(data *model.Good) (*model.Good, error) {
	tx, err := r.storage.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Блокировка
	_, err = tx.Exec(`SELECT * FROM goods WHERE id=$1 AND project_id=$2 FOR UPDATE`, data.Id, data.ProjectId)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(`
			UPDATE goods SET name=$1, description=$2 WHERE id=$3 AND project_id=$4 RETURNING priority, removed, created_at;`,
		data.Name, data.Description, data.Id, data.ProjectId).Scan(&data.Priority, &data.Removed, &data.CreatedAt)
	if err != nil {
		return nil, err
	}

	return data, tx.Commit()
}

func (r *Repo) Delete(id, projectId int) (*model.DeleteResponse, *model.Good, error) {
	var good model.Good
	good.Id = id
	good.ProjectId = projectId

	tx, err := r.storage.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Блокировка
	_, err = tx.Exec(`SELECT * FROM goods WHERE id=$1 AND project_id=$2 FOR UPDATE`, id, projectId)
	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow(`UPDATE goods SET removed=$1 WHERE id=$2 AND project_id=$3 RETURNING name, description, priority, removed, created_at;`,
		true, id, projectId).Scan(&good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	deleteResp := &model.DeleteResponse{Id: id, CampaignId: projectId, Removed: true}

	return deleteResp, &good, tx.Commit()
}

func (r *Repo) UpdatePriority(id, projectId, newPriority int) (*model.ReprioritizeResponse, []model.Good, error) {
	var goods []model.Good

	var good model.Good
	good.Id = id
	good.ProjectId = projectId
	good.Priority = newPriority

	tx, err := r.storage.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Блокировка
	_, err = tx.Exec(`SELECT * FROM goods WHERE id=$1 AND project_id=$2 FOR UPDATE`, id, projectId)
	if err != nil {
		return nil, nil, err
	}

	var priorities []model.Priorities

	// Выполняем запрос на обновление приоритета
	err = tx.QueryRow(`UPDATE goods SET priority=$1 WHERE id=$2 AND project_id=$3 RETURNING name, description, removed, created_at;`,
		newPriority, id, projectId).Scan(&good.Name, &good.Description, &good.Removed, &good.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	priorities = append(priorities, model.Priorities{Id: id, Priority: newPriority})
	goods = append(goods, good)

	// Блокировка
	_, err = tx.Exec(`SELECT * FROM goods WHERE id > $1 FOR UPDATE`, id)
	if err != nil {
		return nil, nil, err
	}

	// Выполняем запрос на обновление приоритета
	_, err = tx.Exec(`UPDATE goods SET priority=$1 WHERE id > $2;`, newPriority, id)
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(`SELECT * FROM goods WHERE id > $1 ORDER BY id;`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// Читаем строки и добавляем их в структуру
	for rows.Next() {
		var priority model.Priorities
		var good model.Good
		err := rows.Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		priority.Id = good.Id
		priority.Priority = good.Priority
		goods = append(goods, good)
		priorities = append(priorities, priority)
	}

	// Проверяем наличие ошибок при чтении строк
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	data := &model.ReprioritizeResponse{Priorities: priorities}

	return data, goods, tx.Commit()
}
