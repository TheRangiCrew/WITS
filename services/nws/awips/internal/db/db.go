package db

import (
	"fmt"
	"log/slog"

	"github.com/surrealdb/surrealdb.go"
)

func New(config DBConfig) (*surrealdb.DB, error) {
	db, err := surrealdb.New(config.Endpoint)
	if err != nil {
		return nil, err
	}

	err = db.Use(config.Namespace, config.Database)
	if err != nil {
		return nil, err
	}

	if config.AsRoot {
		_, err = db.SignIn(&surrealdb.Auth{
			Username: config.Username,
			Password: config.Password,
		})
	} else {
		_, err = db.SignIn(&surrealdb.Auth{
			Namespace: config.Namespace,
			Database:  config.Database,
			Username:  config.Username,
			Password:  config.Password,
		})
	}
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("\033[32mDatabase connected (%s)\033[m", config.Endpoint))

	return db, err
}
