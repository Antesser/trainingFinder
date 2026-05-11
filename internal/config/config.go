package config

import (
	"fmt"
	"log"
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
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Name     string `envconfig:"DB_NAME"`
	SSLMode  string `envconfig:"DB_SSLMODE"`
}

func (c Config) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Name, c.DB.SSLMode)
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env: %w", err)
	}

	portStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("BIP BOP Port should be str: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: (os.Getenv("GRPC_PORT")),
			HTTPPort: (os.Getenv("HTTP_PORT")),
			Host:     (os.Getenv("HOST")),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	return cfg, nil
}
