-- Создание таблицы "projects"
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

INSERT INTO projects (name) VALUES ('Первая запись');

-- Создание таблицы "goods"
CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects (id),
    name TEXT NOT NULL,
    description TEXT,
    priority INT DEFAULT (SELECT COALESCE(MAX(priority), 0) + 1 FROM goods) NOT NULL,
    removed BOOL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_goods_id ON goods (id);
CREATE INDEX IF NOT EXISTS idx_goods_project_id ON goods (project_id);
CREATE INDEX IF NOT EXISTS idx_goods_name ON goods (name);