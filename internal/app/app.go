package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Server.Start(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error(err, "failed to start server")
		}
	}()

	<-done
	app.Logger.Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Stop(ctx); err != nil {
		app.Logger.Error(err, "failed to shutdown server")
	}

	app.Database.Close()

	if err := app.Cache.Close(); err != nil {
		app.Logger.Error(err, "failed to close cache")
	}

	app.Logger.Info("server stopped gracefully")
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
