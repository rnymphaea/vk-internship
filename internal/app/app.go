package app

import (
	"fmt"
	"log"

	"vk-internship/internal/config"
)

func Run() {
	servercfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	storagecfg, err := config.LoadStorageConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(servercfg)
	fmt.Println(storagecfg)
}
