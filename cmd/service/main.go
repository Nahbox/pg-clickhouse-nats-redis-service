package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	msgbRepository "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/msgbroker/nats-streaming"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
	goodsHandler "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/handlers/goods"
	nats "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/infrastructure/nats-streaming"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/infrastructure/postgres"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/infrastructure/redis"
	goodsRepository "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/goods/postgres"
	kvstoreRepository "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/kvstore/redis"
	goodsService "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/service/goods"
)

func main() {
	godotenv.Load()

	cfg, err := config.FromEnv()
	if err != nil {
		log.WithError(err).Fatal("read config from env")
	}

	// Подключение к postgres
	pgdb, err := postgres.New(cfg.PgConfig)
	if err != nil {
		log.WithError(err).Fatal("init postgres")
	}
	defer pgdb.Close()

	// Подключение к redis
	rdb, err := redis.New(cfg.RConfig)
	if err != nil {
		log.WithError(err).Fatal("init redis")
	}
	defer rdb.Close()

	sc, err := nats.New("publisher-client", cfg.NatsConfig.URL())
	if err != nil {
		log.WithError(err).Fatal("init nats publisher")
	}
	defer sc.Close()

	gRepo := goodsRepository.NewGoodsRepo(pgdb)
	rRepo := kvstoreRepository.NewKVStoreRepo(rdb)
	msgbRepo := msgbRepository.NewRepo(sc)

	gService := goodsService.NewService(gRepo, rRepo, msgbRepo)

	gHandler := goodsHandler.NewGoodsHandler(gService)

	router := chi.NewRouter()

	router.Post("/good/create", gHandler.Create)
	router.Patch("/good/update", gHandler.Update)
	router.Delete("/good/remove", gHandler.Remove)
	router.Get("/goods/list", gHandler.GetList)
	router.Patch("/good/reprioritize", gHandler.Reprioritize)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatal("run http server")
		}
	}()
	defer Stop(server)

	log.Infof("started API server on %s", addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	log.Infoln("stopping API server")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("shutdown server")
		os.Exit(1)
	}
}
