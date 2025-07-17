package app

import (
	"log"

	"vk-internship/internal/config"
	"vk-internship/internal/database"
)

type App struct {
	Database database.Database
}

func Run() {
	var app App

	servercfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	storagecfg, err := config.LoadStorageConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.registerDatabase(storagecfg.DBType)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config loaded successfully!")
}
