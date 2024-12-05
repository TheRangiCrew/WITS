package server

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/nwws"
	"github.com/joho/godotenv"
	"github.com/surrealdb/surrealdb.go"
)

type ServerConfig struct {
	DB db.DBConfig
}

type ServerData struct {
	UGC map[string]db.UGC
}

type Server struct {
	DB   *surrealdb.DB
	Data *ServerData
}

func New(config ServerConfig) (*Server, error) {
	godotenv.Load()

	db, err := db.New(config.DB)
	if err != nil {
		return nil, err
	}

	server := Server{
		DB:   db,
		Data: &ServerData{},
	}

	err = server.loadUGC()
	if err != nil {
		return nil, err
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
			h, err := handler.New(server.DB, server.Data.UGC)
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

func (server *Server) loadUGC() error {
	slog.Info("Getting UGC data")

	// Get the latest UGC data
	queryResult, err := surrealdb.Query[[]db.UGC](server.DB, "SELECT * OMIT geometry, centre FROM ugc WHERE valid_to == null", map[string]interface{}{})
	if err != nil {
		return err
	}

	result := *queryResult

	if len(result[0].Result) == 0 {
		return fmt.Errorf("received 0 UGC records")
	}

	data := map[string]db.UGC{}
	for _, ugc := range result[0].Result {
		data[ugc.ID.ID.(string)] = ugc
	}

	server.Data.UGC = data

	slog.Info("Retrieved UGC data")

	return nil
}
