package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"vk-internship/internal/config"
	"vk-internship/internal/database/model"
	"vk-internship/internal/logger"
)

type Redis struct {
	client       *redis.Client
	TTL          time.Duration
	maxFeedItems int
	log          logger.Logger
}

const feedCacheKey = "feed:latest"

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
		client:       client,
		TTL:          cfg.TTL,
		maxFeedItems: cfg.MaxFeedItems,
		log:          log.Component("redis"),
	}

	if err := r.Ping(context.TODO()); err != nil {
		return nil, err
	}

	r.log.Info("connected to redis server")

	return r, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	r.log.Debug("ping redis")
	return r.client.Ping(ctx).Err()
}

func (r *Redis) GetFeed(ctx context.Context) ([]model.Advertisement, error) {
	data, err := r.client.Get(ctx, feedCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			r.log.Debug("feed cache is empty")
			return nil, nil
		}
		r.log.Error(err, "failed to get feed")
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	var ads []model.Advertisement
	if err := json.Unmarshal(data, &ads); err != nil {
		r.log.Error(err, "failed to unmarshal feed")
		return nil, fmt.Errorf("failed to unmarshal feed: %w", err)
	}

	r.log.Debugf("retrieved feed from cache", map[string]interface{}{"count": len(ads)})
	return ads, nil
}

func (r *Redis) SetFeed(ctx context.Context, ads []model.Advertisement) error {
	if len(ads) > r.maxFeedItems {
		ads = ads[:r.maxFeedItems]
	}

	data, err := json.Marshal(ads)
	if err != nil {
		r.log.Warnf("failed to marshal feed", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to marshal feed: %w", err)
	}

	if err := r.client.Set(ctx, feedCacheKey, data, r.TTL).Err(); err != nil {
		r.log.Warnf("failed to set feed cache", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to set feed: %w", err)
	}

	r.log.Debugf("updated feed cache", map[string]interface{}{"count": len(ads)})
	return nil
}

func (r *Redis) UpdateFeed(ctx context.Context, ad model.Advertisement) error {
	ads, err := r.GetFeed(ctx)
	if err != nil {
		return err
	}

	updatedAds := make([]model.Advertisement, 0, r.maxFeedItems)
	updatedAds = append(updatedAds, ad)

	if len(ads) >= r.maxFeedItems {
		ads = ads[:r.maxFeedItems-1]
	}
	updatedAds = append(updatedAds, ads...)

	return r.SetFeed(ctx, updatedAds)
}

func (r *Redis) InvalidateFeed(ctx context.Context) error {
	if err := r.client.Del(ctx, feedCacheKey).Err(); err != nil {
		r.log.Warnf("failed to invalidate feed cache", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to invalidate feed: %w", err)
	}
	r.log.Debug("invalidated feed cache")
	return nil
}

func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		r.log.Error(err, "failed to close connection")
		return err
	}
	r.log.Info("redis connection closed successfully")
	return nil
}

func (r *Redis) GetMaxFeedItems() int {
	return r.maxFeedItems
}
