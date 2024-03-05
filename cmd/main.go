package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/go-chi/chi/v5"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/config"
	"github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/db/postgres"
	goodsHandler "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/handlers/goods"
	goodsRepository "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/repository/goods/postgres"
	goodsService "github.com/Nahbox/pg-clickhouse-nats-redis-service/internal/service/goods"
)

func main() {
	if os.Getenv("APP_ENV") == "local" {
		godotenv.Load()
	}

	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal("read config from env", err)
	}

	database, err := postgres.New(cfg.PgConfig)
	if err != nil {
		log.Fatal("init db", err)
	}
	defer database.Close()

	gRepo, err := goodsRepository.NewGoodsRepo(database)
	if err != nil {
		log.Fatal(err)
	}

	gService, err := goodsService.NewService(gRepo)
	if err != nil {
		log.Fatal(err)
	}

	gHandler, err := goodsHandler.NewGoodsHandler(gService)
	if err != nil {
		log.Fatal(err)
	}

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
