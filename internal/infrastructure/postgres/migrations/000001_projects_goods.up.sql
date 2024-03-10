-- Создание таблицы "projects"
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

INSERT INTO projects (name) VALUES ('First record');
INSERT INTO projects (name) VALUES ('Second record');
INSERT INTO projects (name) VALUES ('Third record');
INSERT INTO projects (name) VALUES ('Fourth record');
INSERT INTO projects (name) VALUES ('Fifth record');

-- Создание таблицы "goods"
CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects (id),
    name TEXT NOT NULL,
    description TEXT,
    priority SERIAL NOT NULL,
    removed BOOL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_goods_id_project_id ON goods (id, project_id);
CREATE INDEX IF NOT EXISTS idx_goods_name ON goods (name);