package config

import "os"

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	RestPort string
	GRPCPort string
}

type JWTConfig struct {
	Secret string
	Issuer string
}

type DatabaseConfig struct {
	URL string
}

// Message Queue Config
type RabbitMQConfig struct {
	URL        string
	Queue      string
	Exchange   string
	RoutingKey string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			RestPort: getEnv("SERVER_PORT", "8080"),
			GRPCPort: getEnv("GRPC_PORT", "50052"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET_ACCESS_TOKEN", "secret-jwt"),
			Issuer: getEnv("JWT_ISSUER", "sidingin-app"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sidingin_rating_service?sslmode=disable"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Queue:      getEnv("RABBITMQ_QUEUE", "rating_service.order_completed.queue"),
			Exchange:   getEnv("RABBITMQ_EXCHANGE", ""),
			RoutingKey: getEnv("RABBITMQ_ROUTING_KEY", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
