package server

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/TheRangiCrew/WITS/services/parsing/awips/internal/handler"
)

func ParseText(filename string, minLog int) {
	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	text := string(data)

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

	slog.Info(fmt.Sprintf("Parsing %s", filename))

	h.Handle(text, time.Now())

}
