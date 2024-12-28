package server

import (
	"log/slog"
	"os"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/nwws"
	"github.com/surrealdb/surrealdb.go"
)

type ServerConfig struct {
	DB     db.DBConfig
	MinLog int
}

type Server struct {
	DB     *surrealdb.DB
	MinLog int
}

func New(config ServerConfig) (*Server, error) {
	db, err := db.New(config.DB)
	if err != nil {
		return nil, err
	}

	server := Server{
		DB:     db,
		MinLog: config.MinLog,
	}

	return &server, nil
}

func NWWS(config ServerConfig) {
	server, err := New(config)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	nwwsoi, err := nwws.New(&nwws.Config{
		Server:   os.Getenv("NWWSOI_Server") + ":5222",
		Room:     os.Getenv("NWWSOI_Room"),
		User:     os.Getenv("NWWSOI_User"),
		Pass:     os.Getenv("NWWSOI_Pass"),
		Resource: os.Getenv("NWWSOI_Resource"),
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	queue := make(chan *nwws.Message)
	errChan := make(chan error)

	go nwwsoi.Start(queue)

	go func() {
		for message := range queue {
			h, err := handler.New(server.DB, server.MinLog)
			if err != nil {
				errChan <- err
				return
			}
			go h.Handle(message.Text, message.ReceivedAt.UTC())
		}
	}()

	for err := range errChan {
		slog.Error(err.Error())
	}
}
