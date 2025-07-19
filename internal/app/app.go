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

	err = app.registerLogger(loggercfg)
	if err != nil {
		log.Fatal(err)
	}

	servercfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.registerServer(servercfg, app.Logger)

	storagecfg, err := config.LoadStorageConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.registerDatabase(storagecfg.DBType, app.Logger)
	if err != nil {
		log.Fatal(err)
	}

	err = app.registerCache(storagecfg.CacheType, app.Logger)
	if err != nil {
		log.Fatal(err)
	}

	err = app.registerLogger(loggercfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("config loaded successfully")

	app.Server.Start()
}
