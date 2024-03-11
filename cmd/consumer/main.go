package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/infrastructure/clickhouse"
	nats "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/infrastructure/nats-streaming"
	logsRepo "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/logs/clickhouse"
	msgBrokerRepo "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/msgbroker/nats-streaming"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/service/logs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	godotenv.Load()

	cfg, err := config.FromEnv()
	if err != nil {
		log.WithError(err).Fatal("read config from env")
	}

	chConn, err := clickhouse.New(cfg.CHConfig)
	if err != nil {
		log.WithError(err).Fatal("init clickhouse")
	}
	defer chConn.Close()

	natsConn, err := nats.New("subscriber-client", cfg.NatsConfig.URL())
	if err != nil {
		log.WithError(err).Fatal("init nats subscriber")
	}
	defer natsConn.Close()

	logsRepo := logsRepo.NewRepo(chConn)
	msgBrokerRepo := msgBrokerRepo.NewRepo(natsConn)

	service := logs.NewService(logsRepo, msgBrokerRepo)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := service.SaveLogs(ctx)
		if err != nil {
			log.WithError(err).Fatal("save logs")
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	cancel()
	wg.Wait()

	log.Infoln("stopping consumer")
}
