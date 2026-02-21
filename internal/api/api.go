package api

import (
	"context"
	"device-parser-logs/internal/config"
	"device-parser-logs/internal/controller"
	"device-parser-logs/internal/generator"
	"device-parser-logs/internal/parser"
	"device-parser-logs/internal/repository"
	"device-parser-logs/internal/service"
	"device-parser-logs/internal/watcher"
	"device-parser-logs/pkg/client"
	"device-parser-logs/producer"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Api struct {
	logger     *slog.Logger
	service    *service.Service
	watcher    *watcher.Pool
	dbClient   *mongo.Client
	httpServer *controller.Server
	producer   *producer.Producer
}

func New() *Api {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to init configs %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	clientMongo, err := client.NewClient(ctx, cfg.Host, cfg.Port, cfg.User, cfg.Password)
	if err != nil {
		log.Fatalf("failed to connect mongoDB %v", err)
	}

	repo := repository.New(ctx, cfg.Name, clientMongo)

	parser := parser.New(cfg.DirectoryTsv)

	generator := generator.New(cfg.DirectoryPdf)

	producer := producer.New(cfg.DirectoryTsv)

	service := service.New(parser, repo, generator, logger)

	watcher := watcher.NewPool(ctx, cfg.Workers, cfg.DirectoryTsv, cfg.Interval, logger, service)

	handlers := controller.NewHandler(service, logger)
	httpServer := controller.NewServer(handlers, cfg.Addr)

	return &Api{
		logger:     logger,
		service:    service,
		watcher:    watcher,
		dbClient:   clientMongo,
		httpServer: httpServer,
		producer:   producer,
	}
}

func (a *Api) RunApi() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a.watcher.RunPool()

	a.producer.Produce()

	go func() {
		a.logger.Info("starting http server")
		err := a.httpServer.StartHttpServer()
		if err != nil {
			log.Fatalf("failed to start http server %v", err)
		}
	}()

	<-sigChan

	a.watcher.Stop()

	err := client.CloseMongoDB(a.dbClient)
	if err != nil {
		a.logger.Error("failed to disconnect to mongoDB", "error", err.Error())
	}
}
