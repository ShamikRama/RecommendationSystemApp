package app

import (
	"User_Service/internal/api"
	"User_Service/internal/config"
	"User_Service/internal/logger"
	"User_Service/internal/repository"
	"User_Service/internal/service"
	"User_Service/internal/storage"
	"User_Service/pkg/kafka"
	"User_Service/pkg/product"
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
		Addr:         cfg.Address,
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

	err = kafka.CreateTopic("kafka:9092")
	if err != nil {
		log.Fatalf("failed create topic: %v", err)
	}

	// Инициализация логера
	logs := logger.NewLogger()
	logs.Info("Logger initialized")

	// Инициализация кафка продюсера
	kafkaProducer := kafka.NewKafka(*logs)
	logs.Info("Kafka producer initialized")

	//Инициализация репозитория
	repos := repository.NewRepository(db, logs)
	logs.Info("Repository initialized")

	// Инициализация сервисов
	services := service.NewService(*repos, logs, kafkaProducer)
	logs.Info("Services initialized")

	// Инициализация клиента продуктов
	productClient := product.NewClient(*logs) // Создаем клиент через интерфейс
	logs.Info("Product client initialized")

	// Инициализация обработчиков
	handlers := api.NewApi(services, productClient, *logs)
	logs.Info("Handlers initialized")

	// Инициализация маршрутов
	router := handlers.InitRoutes()
	logs.Info("Routes initialized")

	// Инициализация HTTP-сервера
	server := InitHttpServer(*cfg, router)
	logs.Info("Server initialized")

	// Запуск HTTP-сервера в горутине
	go func() {
		logs.Info("Server staring")

		if err := RunServer(server); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	kafkaProducer.Close()
	logs.Info("Kafka producer closed")

	db.Close()
	logs.Info("Database closed")
}
