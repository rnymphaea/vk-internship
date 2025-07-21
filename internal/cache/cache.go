package cache

import (
	"context"

	"vk-internship/internal/database/model"
)

type Cache interface {
	Ping(ctx context.Context) error
	GetFeed(ctx context.Context) ([]model.Advertisement, error)
	SetFeed(ctx context.Context, ads []model.Advertisement) error
	UpdateFeed(ctx context.Context, ad model.Advertisement) error
	InvalidateFeed(ctx context.Context) error
	Close() error
}
