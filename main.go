package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
}

func (c Config) PostgresURL() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("error loading .env: %w", err)
    }

    port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     port,
        DBUser:     getEnv("DB_USER", "admin"),
        DBPassword: getEnv("DB_PASSWORD", "falsePass"),
        DBName:     getEnv("DB_NAME", "godb"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func main() {
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    dbURL := cfg.PostgresURL()

    conn, err := pgx.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatal("Unable to connect:", err)
    }
    defer conn.Close(context.Background())

    fmt.Println("Successfully connected to", cfg.DBName)
}
