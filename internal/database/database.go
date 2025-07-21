package database

import (
	"context"
	"errors"

	"vk-internship/internal/database/model"
)

type Database interface {
	Ping(ctx context.Context) error

	CreateUser(user *model.User) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)

	CreateAd(ad *model.Advertisement) (*model.Advertisement, error)
	GetAds(ctx context.Context, sortBy, order string, minPrice, maxPrice *int, page, pageSize int) ([]*model.Advertisement, int, error)
	GetAd(ctx context.Context, id string) (*model.Advertisement, error)
	UpdateAd(ctx context.Context, ad *model.Advertisement) (*model.Advertisement, error)
	DeleteAd(ctx context.Context, id, authorID string) error

	Close()
}

var (
	ErrUserExists                 = errors.New("username already exists")
	ErrUserNotFound               = errors.New("user not found")
	ErrAdNotFound                 = errors.New("advertisement not found")
	ErrAdNotFoundOrNotOwnedByUser = errors.New("advertisement not found or not owned by user")
)
