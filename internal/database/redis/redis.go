package redis

import (
	"LoveMusic/internal/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(cfg *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return client, nil

}
