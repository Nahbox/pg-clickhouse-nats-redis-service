nats-serv:
	docker compose up nats-streaming -d

postgres:
	docker compose up postgres -d

rd:
	docker compose up redis -d

ch:
	docker compose up clickhouse -d
	sleep 4

infra: postgres rd nats-serv ch

service:
	go run cmd/consumer/main.go & go run cmd/service/main.go

all: infra service

tmp:
	go run cmd/consumer/main.go &

clean:
	docker compose down &
	pkill -f "go run cmd/consumer/main.go" &
	pkill -f "go run cmd/service/main.go" &
