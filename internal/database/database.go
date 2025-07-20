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
}

var (
	ErrUserExists   = errors.New("username already exists")
	ErrUserNotFound = errors.New("user not found")
)
