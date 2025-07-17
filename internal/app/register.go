package app

import (
	"fmt"

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
