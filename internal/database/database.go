package database

import (
	"context"
	"errors"

	"vk-internship/internal/database/model"
)

type Database interface {
	Ping(ctx context.Context) error
	CreateUser(user *model.User) (*model.User, error)
}

var (
	ErrUserExists = errors.New("username already exists")
)
