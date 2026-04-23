package storage

import (
	"context"
	"fmt"

	"server/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, cfg config.Config) (*redis.Client, error) {
	addr := cfg.Redis.Addr()
	if addr == "" {
		return nil, fmt.Errorf("redis addr is empty")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}
