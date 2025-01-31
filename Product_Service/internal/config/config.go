package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer    `env:"HTTP_SERVER"`
	Database      `env:"DATABASE"`
	PgcConnString string `env:"PG_DSN"`
	MigConnString string `env:"MIGRATION_DSN"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS"`
	Topic   string   `env:"KAFKA_TOPIC"`
	GroupID string   `env:"KAFKA_GROUP_ID"`
}

// HTTPServer содержит настройки HTTP-сервера
type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"0.0.0.0:8082"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

// Database содержит настройки базы данных
type Database struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DBName   string `env:"POSTGRES_DBNAME"`
	DSN      string `env:"PG_DSN"`
}

func MustLoad() *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH environment variable is not set")
	}

	err := godotenv.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read environment variables: %v", err)
	}

	return &cfg
}
