package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/Antesser/trainingFinder/internal/config"
	redisModel "github.com/Antesser/trainingFinder/internal/model/redis"
	telegramModel "github.com/Antesser/trainingFinder/internal/model/telegram"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis(ctx context.Context, cfg config.ServerConfig) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
}

func SessionKey(chatID int64) string {
	return fmt.Sprintf("%s%d", redisModel.KeyPrefix, chatID)
}

func LoadSession(ctx context.Context, chatID int64) (*redisModel.Session, error) {
	key := SessionKey(chatID)
	data, err := Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &redisModel.Session{State: telegramModel.StateIdle}, nil
	}
	return &redisModel.Session{
		State:    data["state"],
		Login:    data["login"],
		Password: data["password"],
		Role:     data["role"],
	}, nil
}

func SaveSession(ctx context.Context, chatID int64, s *redisModel.Session) error {
	key := SessionKey(chatID)
	return Rdb.HSet(ctx, key,
		"state", s.State,
		"login", s.Login,
		"password", s.Password,
		"role", s.Role,
	).Err()
}
