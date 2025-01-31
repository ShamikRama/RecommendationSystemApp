package app

import (
	"Analytics_Service/internal/api"
	"Analytics_Service/internal/config"
	"Analytics_Service/internal/logger"
	"Analytics_Service/internal/repository"
	"Analytics_Service/internal/service"
	"Analytics_Service/internal/storage"
	"Analytics_Service/pkg/kafka"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
)

type HttpServer struct {
	server http.Server
}

func InitHttpServer(cfg *config.Config, router http.Handler) *HttpServer {
	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &HttpServer{
		server: srv,
	}
}

func RunServer(server *HttpServer) error {
	if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func RunApp() {
	cfg := config.MustLoad()

	// Подключение к ClickHouseсду
	db, err := storage.InitClickhouse(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db(postgres): %v", err)
	}

	// Инициализация логера
	logs := logger.NewLogger()
	logs.Info("Logger initialized")

	// Инициализация репозиториев
	clientRepo := repository.NewClientRepo(db, logs)
	analRepo := repository.NewAnalRepository(db, logs)
	logs.Info("Repositories initialized")

	// Инициализация сервисов
	clientService := service.NewClientService(clientRepo, *logs)
	handlerService := service.NewHandler(*logs, analRepo)
	logs.Info("Services initialized")

	// Инициализация Kafka Consumer
	kafkaConsumer := kafka.NewKafkaConsumer(
		*logs,
		&handlerService,
		cfg.KafkaConfig.Brokers,
		cfg.KafkaConfig.Topics,
		cfg.KafkaConfig.GroupID,
	)
	logs.Info("Kafka consumer initialized")

	go func() {
		logs.Info("Starting Kafka consumer")
		kafkaConsumer.ReadMessages()
	}()

	// Инициализация API
	apiHandlers := api.NewApi(*clientService)
	logs.Info("API handlers initialized")

	// Инициализация маршрутов
	router := api.InitRoutes(apiHandlers)
	logs.Info("Routes initialized")

	// Инициализация HTTP-сервера
	server := InitHttpServer(cfg, router)
	logs.Info("Server initialized")

	go func() {
		logs.Info("Server starting")
		if err := RunServer(server); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logs.Info("Shutting down server...")

	kafkaConsumer.Close()
	logs.Info("Closed Kafka consumer")

	db.Close()
	logs.Info("Database closed")
}
