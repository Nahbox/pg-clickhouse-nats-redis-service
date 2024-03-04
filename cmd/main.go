package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	conn, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?username=default&password=")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Создаем таблицу
	_, err = conn.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS logs (
			id UInt32,
			message String
		) ENGINE = MergeTree()
		ORDER BY id
	`)
	if err != nil {
		log.Fatal(err)
	}
}
