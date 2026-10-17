package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/Antesser/trainingFinder/internal/config"
	telegramModel "github.com/Antesser/trainingFinder/internal/model/telegram"
	"github.com/redis/go-redis/v9"
)

const (
	keyPrefix = "bot:session:"
)

type cache struct {
	Rdb *redis.Client
}

func InitRedis(ctx context.Context, cfg config.ServerConfig) *cache {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	return &cache{Rdb: Rdb}
}

func sessionKey(chatID int64) string {
	return fmt.Sprintf("%s%d", keyPrefix, chatID)
}

func (r *cache) LoadSession(ctx context.Context, chatID int64) (*telegramModel.Session, error) {
	key := sessionKey(chatID)
	data, err := r.Rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &telegramModel.Session{State: telegramModel.StateIdle}, nil
	}
	return &telegramModel.Session{
		State:    data["state"],
		Login:    data["login"],
		Password: data["password"],
		Role:     data["role"],
	}, nil
}

func (r *cache) SaveSession(ctx context.Context, chatID int64, s *telegramModel.Session) error {
	key := sessionKey(chatID)
	return r.Rdb.HSet(ctx, key,
		"state", s.State,
		"login", s.Login,
		"password", s.Password,
		"role", s.Role,
	).Err()
}
func (r *cache) DelSession(ctx context.Context, chatID int64) error {
	return r.Rdb.Del(ctx, sessionKey(chatID)).Err()
}
