package app

import (
	"log"

	"vk-internship/internal/cache"
	"vk-internship/internal/config"
	"vk-internship/internal/database"
	"vk-internship/internal/logger"
)

type App struct {
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

	_, err = config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

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
}
