package queue

import (
	"context"
	"go-user-service/config"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient( cfg config.RedisConfig) (*redis.Client, error) {
	ctx := context.Background()

	rdb  := redis.NewClient(&redis.Options{
		Addr: cfg.Host+":"+cfg.Port,
		Password: cfg.Password,
	})

	err := rdb.Set(ctx, "key", "Value", 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
	
    return rdb, nil
}