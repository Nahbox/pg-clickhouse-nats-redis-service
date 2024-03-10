nats-serv:
	docker compose up nats-streaming -d
	sleep 5

postgres:
	docker compose up postgres -d
	sleep 5

rd:
	docker compose up redis -d
	sleep 5

ch:
	docker compose up clickhouse -d
	sleep 5

service:
	go run cmd/service/main.go

infrastructure: postgres rd nats-serv ch

all: postgres rd nats-serv ch service

clean:
	docker compose down