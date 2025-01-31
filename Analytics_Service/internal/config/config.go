package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer           `env:"HTTP_SERVER"`
	KafkaConfig          `env:"KAFKA"`
	ConnStringCLickHouse string `env:"CH_DSN"`
	ConnStrPSQL          string `env:"PG_DSN"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS"`
	Topics  []string `env:"KAFKA_TOPICS"`
	GroupID string   `env:"KAFKA_GROUP_ID"`
}

// HTTPServer содержит настройки HTTP-сервера
type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"0.0.0.0:8085"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

// ClickHouse содержит настройки ClickHouse
type ClickHouse struct {
	Addr     string `env:"CLICKHOUSE_ADDR" env-default:"clickhouse:9000"`
	DB       string `env:"CLICKHOUSE_DB" env-default:"analytics_db"`
	User     string `env:"CLICKHOUSE_USER" env-default:"admin"`
	Password string `env:"CLICKHOUSE_PASSWORD" env-default:"password"`
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

	// Чтение переменных окружения в структуру Config
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Failed to read environment variables: %v", err)
	}

	cfg.KafkaConfig.Topics = strings.Split(os.Getenv("KAFKA_TOPICS"), ",")

	return &cfg
}
