package server

import (
	"log/slog"
	"os"

	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/db"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/handler"
	"github.com/TheRangiCrew/WITS/services/nws/awips/internal/nwws"
)

type ServerConfig struct {
	DB db.DBConfig
}

func Start(config ServerConfig) {

	nwwsoi, err := nwws.New(&nwws.Config{
		Server:   os.Getenv("NWWSOI_Server") + ":5222",
		Room:     os.Getenv("NWWSOI_Room"),
		User:     os.Getenv("NWWSOI_User"),
		Pass:     os.Getenv("NWWSOI_Pass"),
		Resource: os.Getenv("NWWSOS_Resource"),
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	handler, err := handler.New(config.DB)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	queue := make(chan *nwws.Message)

	go nwwsoi.Start(queue)

	go func() {
		for message := range queue {
			go handler.Handle(message.Text, message.ReceivedAt.UTC())
		}
	}()

	for err := range handler.ErrChan {
		slog.Error(err.Error())
	}
}
