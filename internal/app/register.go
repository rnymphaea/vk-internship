package app

import (
	"fmt"

	"vk-internship/internal/cache"
	"vk-internship/internal/cache/redis"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/database/postgres"
	"vk-internship/internal/logger"
	zerologger "vk-internship/internal/logger/zerolog"
	"vk-internship/internal/server"
)

func (app *App) registerDatabase(dbType string, log logger.Logger) error {
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

		db, err = postgres.New(cfg, log)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("database type [%s] is not supported", dbType)
	}

	app.Database = db
	return err
}

func (app *App) registerCache(cacheType string, log logger.Logger) error {
	var (
		cache cache.Cache
		err   error
	)

	switch cacheType {
	case "redis":
		cfg, err := config.LoadRedisConfig()
		if err != nil {
			return err
		}

		cache, err = redis.New(cfg, log)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("cache type [%s] is not supported", cacheType)
	}

	app.Cache = cache
	return err
}

func (app *App) registerLogger(cfg *config.LoggerConfig) error {
	var (
		logger logger.Logger
		err    error
	)

	switch cfg.Type {
	case "zerolog":
		logger = zerologger.New(cfg)
	default:
		return fmt.Errorf("logger type [%s] is not supported", cfg.Type)
	}

	app.Logger = logger
	return err
}

func (app *App) registerServer(servercfg *config.ServerConfig, log logger.Logger) {
	srv := server.New(servercfg, nil, log)
	app.Server = srv
}
