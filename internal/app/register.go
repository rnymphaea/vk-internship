package app

import (
	"fmt"

	"vk-internship/internal/cache"
	"vk-internship/internal/cache/redis"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/database/postgres"
)

func (app *App) registerDatabase(dbType string) error {
	var (
		db  database.Database
		err error
	)

	switch dbType {
	case "postgres":
		cfg, err := config.LoadPostgresConfig()
		if err != nil {
			return err
		}

		db, err = postgres.New(cfg)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("database type [%s] is not supported", dbType)
	}

	app.Database = db
	return err
}

func (app *App) registerCache(cacheType string) error {
	var (
		cache cache.Cache
		err   error
	)

	switch cacheType {
	case "redis":
		cfg, err := config.LoadRedisConfig()
		fmt.Println(cfg)
		if err != nil {
			return err
		}

		cache, err = redis.New(cfg)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("cache type [%s] is not supported", cacheType)
	}

	app.Cache = cache
	return err
}
