package server

import (
	"log/slog"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/handler"
)

func ParseText(text string, minLog int) {

	config := ServerConfig{
		MinLog: minLog,
	}

	server, err := New(config)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	h, err := handler.New(server.DB, server.MinLog)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	h.Handle(text, time.Now())

}
