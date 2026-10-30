package config

import "os"

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
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

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			RestPort: getEnv("SERVER_PORT", "8080"),
			GRPCPort: getEnv("GRPC_PORT", "50051"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET_ACCESS_TOKEN", "secret-jwt"),
			Issuer: getEnv("JWT_ISSUER", "sidingin-app"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sidingin_rating_service?sslmode=disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
