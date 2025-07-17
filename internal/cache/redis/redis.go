package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"vk-internship/internal/config"
)

type Redis struct {
	client *redis.Client
}

func New(cfg *config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	r := &Redis{
		client: client,
	}

	if err := r.Ping(context.Background()); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return r, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
