package queue

import (
	"context"
	"go-user-service/config"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient( cfg config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	ctx := context.Background()

	rdb  := redis.NewClient(&redis.Options{
		Addr: cfg.Host+":"+cfg.Port,
		Password: cfg.Password,
	})

	pong, _ := rdb.Ping(ctx).Result()
    if pong != "" {
		logger.Info("Redis is up")
	}
	
    return rdb, nil
}