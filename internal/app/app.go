package app

import (
	"log"

	"vk-internship/internal/cache"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
	"vk-internship/internal/server"
)

type App struct {
	Server   *server.Server
	Database database.Database
	Cache    cache.Cache
	Logger   logger.Logger
}

func Run() {
	var app App

	loggercfg, err := config.LoadLoggerConfig()
	if err != nil {
		log.Fatal(err)
	}

	servercfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	storagecfg, err := config.LoadStorageConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.registerComponents(loggercfg, servercfg, storagecfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("config loaded successfully")
	app.Server.Start()
}

func (app *App) registerComponents(loggercfg *config.LoggerConfig, servercfg *config.ServerConfig, storagecfg *config.StorageConfig) error {
	err := app.registerLogger(loggercfg)
	if err != nil {
		return err
	}

	err = app.registerDatabase(storagecfg.DBType, app.Logger)
	if err != nil {
		return err
	}

	err = app.registerCache(storagecfg.CacheType, app.Logger)
	if err != nil {
		return err
	}

	app.registerServer(servercfg, app.Logger)

	return nil
}
