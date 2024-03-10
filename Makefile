nats-serv:
	docker compose up nats-streaming -d
	sleep 4

postgres:
	docker compose up postgres -d
	sleep 4

rd:
	docker compose up redis -d
	sleep 4

ch:
	docker compose up clickhouse -d
	sleep 4

consumer:
	go run cmd/consumer/main.go &

service:
	go run cmd/service/main.go &

all: postgres rd nats-serv ch service consumer

clean:
	docker compose down