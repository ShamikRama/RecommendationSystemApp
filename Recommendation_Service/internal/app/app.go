package app

import (
	"Recommendation_Service/internal/api"
	"Recommendation_Service/internal/config"
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/repository"
	"Recommendation_Service/internal/service"
	"Recommendation_Service/internal/storage"
	"Recommendation_Service/pkg/kafka"
	"Recommendation_Service/pkg/redis"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type HttpServer struct {
	server http.Server
}

func InitHttpServer(cfg config.Config, router http.Handler) *HttpServer {
	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout * time.Second,
		WriteTimeout: cfg.HTTPServer.Timeout * time.Second,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout * time.Second,
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

	// Подключение к бд
	db, err := storage.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Инициализация логера
	logs := logger.NewLogger()
	logs.Info("Logger initialized")

	// Инициализация репозитория
	repos := repository.NewProductRepository(db, logs)
	logs.Info("Repository initialized")

	redisClient := redis.NewRedisClient()
	logs.Info("Redis initialized")

	// Инициализация сервисов
	services := service.NewRecService(repos, logs, redisClient)
	logs.Info("Services initialized")

	// Инициализация Kafka Consumer
	kafkaHandler := service.NewHandler(*logs, repos, redisClient)
	kafkaConsumer := kafka.NewKafkaConsumer(
		*logs,
		&kafkaHandler,
		cfg.KafkaConfig.Brokers,
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.GroupID,
	)
	logs.Info("Kafka consumer initialized")

	go func() {
		logs.Info("Starting Kafka consumer")
		kafkaConsumer.ReadMessages()
	}()

	// Инициализация обработчиков
	handlers := api.NewRecApi(services, *logs)
	logs.Info("Handlers initialized")

	// Инициализация маршрутов
	router := handlers.InitRoutes()
	logs.Info("Routes initialized")

	// Инициализация HTTP-сервера
	server := InitHttpServer(*cfg, router)
	logs.Info("Server initialized")

	// Запуск HTTP-сервера в горутине
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
	logs.Info("Kafka consumer closed")

	redisClient.Close()
	logs.Info("Redis client closed")

	db.Close()
	logs.Info("Database closed")

}
