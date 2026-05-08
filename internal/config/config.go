package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	GRPCPort string `envconfig:"GRPC_PORT"`
	HTTPPort string `envconfig:"HTTP_PORT"`
	Host     string `envconfig:"HOST"`
}

type DBConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string ``
	Name     string
	SSLMode  string
}

func (c Config) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: (getEnv("GRPC_PORT", ":9090")),
			HTTPPort: (getEnv("HTTP_PORT", ":8080")),
			Host:     (getEnv("HOST", "localhost")),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getIntEnv("DB_PORT", 5432),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "godb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func getIntEnv(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}
