package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"vk-internship/internal/config"
	"vk-internship/internal/logger"
)

type Redis struct {
	client *redis.Client
	log    logger.Logger
}

func New(cfg *config.RedisConfig, log logger.Logger) (*Redis, error) {
	log.Debug("creating new redis client")

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
		log:    log.Component("redis"),
	}

	if err := r.Ping(context.Background()); err != nil {
		return nil, err
	}

	r.log.Info("connected to redis server")

	return r, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	r.log.Debug("ping redis")
	return r.client.Ping(ctx).Err()
}
