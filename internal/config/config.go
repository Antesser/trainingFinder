package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig

	TrainingOutboxProcessEnabled  bool          `envconfig:"TRAINING_OUTBOX_PROCESS_ENABLED" default:"false"`
	TrainingOutboxProcessDuration time.Duration `envconfig:"TRAINING_OUTBOX_PROCESS_DURATION"` // 5s 2h
}

type ServerConfig struct {
	GRPCPort            string        `envconfig:"GRPC_PORT"`
	HTTPPort            string        `envconfig:"HTTP_PORT"`
	Host                string        `envconfig:"HOST"`
	Secret              string        `envconfig:"SECRET"`
	AccessTokenDuration time.Duration `envconfig:"ACCESS_TOKEN_DURATION"`
	AuthConfigPath      string        `envconfig:"AUTH_CONFIG_PATH"`
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
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to process envconfig: %w", err)
	}
	return &cfg, nil
}
